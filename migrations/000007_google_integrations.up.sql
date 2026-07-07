-- =========================================================
-- INTEGRATION SOURCES
-- =========================================================

CREATE TABLE integration_sources (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                     code VARCHAR(100) NOT NULL,
                                     name VARCHAR(150) NOT NULL,

                                     provider VARCHAR(50) NOT NULL,
                                     source_type VARCHAR(50) NOT NULL,

                                     external_id VARCHAR(255),

                                     status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',

                                     config JSONB NOT NULL DEFAULT '{}'::JSONB,

                                     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                     CONSTRAINT uq_integration_sources_code
                                         UNIQUE (code),

                                     CONSTRAINT chk_integration_sources_code_not_blank
                                         CHECK (BTRIM(code) <> ''),

                                     CONSTRAINT chk_integration_sources_name_not_blank
                                         CHECK (BTRIM(name) <> ''),

                                     CONSTRAINT chk_integration_sources_provider
                                         CHECK (
                                             provider IN (
                                                          'GOOGLE_APPS_SCRIPT',
                                                          'GOOGLE_SHEETS',
                                                          'GOOGLE_FORMS',
                                                          'GOOGLE_DRIVE'
                                                 )
                                             ),

                                     CONSTRAINT chk_integration_sources_source_type_not_blank
                                         CHECK (BTRIM(source_type) <> ''),

                                     CONSTRAINT chk_integration_sources_status
                                         CHECK (
                                             status IN (
                                                        'ACTIVE',
                                                        'INACTIVE'
                                                 )
                                             )
);

CREATE INDEX idx_integration_sources_provider
    ON integration_sources (provider);

CREATE INDEX idx_integration_sources_status
    ON integration_sources (status);

-- =========================================================
-- INTEGRATION SYNC RUNS
-- =========================================================

CREATE TABLE integration_sync_runs (
                                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                       source_id UUID NOT NULL,

                                       direction VARCHAR(20) NOT NULL,
                                       status VARCHAR(30) NOT NULL DEFAULT 'PENDING',

                                       started_at TIMESTAMPTZ,
                                       completed_at TIMESTAMPTZ,

                                       records_received INTEGER NOT NULL DEFAULT 0,
                                       records_processed INTEGER NOT NULL DEFAULT 0,
                                       records_failed INTEGER NOT NULL DEFAULT 0,

                                       error_message TEXT,

                                       metadata JSONB NOT NULL DEFAULT '{}'::JSONB,

                                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                       CONSTRAINT fk_integration_sync_runs_source
                                           FOREIGN KEY (source_id)
                                               REFERENCES integration_sources(id)
                                               ON DELETE CASCADE,

                                       CONSTRAINT chk_integration_sync_runs_direction
                                           CHECK (
                                               direction IN (
                                                             'INBOUND',
                                                             'OUTBOUND'
                                                   )
                                               ),

                                       CONSTRAINT chk_integration_sync_runs_status
                                           CHECK (
                                               status IN (
                                                          'PENDING',
                                                          'RUNNING',
                                                          'SUCCESS',
                                                          'PARTIAL',
                                                          'FAILED'
                                                   )
                                               ),

                                       CONSTRAINT chk_integration_sync_runs_counts
                                           CHECK (
                                               records_received >= 0
                                                   AND records_processed >= 0
                                                   AND records_failed >= 0
                                               ),

                                       CONSTRAINT chk_integration_sync_runs_date_range
                                           CHECK (
                                               completed_at IS NULL
                                                   OR started_at IS NULL
                                                   OR completed_at >= started_at
                                               )
);

CREATE INDEX idx_integration_sync_runs_source_id
    ON integration_sync_runs (source_id);

CREATE INDEX idx_integration_sync_runs_status
    ON integration_sync_runs (status);

CREATE INDEX idx_integration_sync_runs_created_at
    ON integration_sync_runs (created_at DESC);

-- =========================================================
-- INTEGRATION EVENTS
-- =========================================================

CREATE TABLE integration_events (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                    source_id UUID NOT NULL,
                                    sync_run_id UUID,

                                    external_event_id VARCHAR(255),

                                    event_type VARCHAR(100) NOT NULL,

                                    payload JSONB NOT NULL,
                                    payload_hash BYTEA,

                                    status VARCHAR(30) NOT NULL DEFAULT 'RECEIVED',

                                    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                    processed_at TIMESTAMPTZ,

                                    error_message TEXT,

                                    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                    CONSTRAINT fk_integration_events_source
                                        FOREIGN KEY (source_id)
                                            REFERENCES integration_sources(id)
                                            ON DELETE CASCADE,

                                    CONSTRAINT fk_integration_events_sync_run
                                        FOREIGN KEY (sync_run_id)
                                            REFERENCES integration_sync_runs(id)
                                            ON DELETE SET NULL,

                                    CONSTRAINT chk_integration_events_event_type_not_blank
                                        CHECK (BTRIM(event_type) <> ''),

                                    CONSTRAINT chk_integration_events_status
                                        CHECK (
                                            status IN (
                                                       'RECEIVED',
                                                       'PROCESSING',
                                                       'PROCESSED',
                                                       'FAILED',
                                                       'IGNORED'
                                                )
                                            ),

                                    CONSTRAINT chk_integration_events_processed_at
                                        CHECK (
                                            processed_at IS NULL
                                                OR processed_at >= received_at
                                            )
);

CREATE UNIQUE INDEX uq_integration_events_external_event
    ON integration_events (
                           source_id,
                           external_event_id
        )
    WHERE external_event_id IS NOT NULL;

CREATE INDEX idx_integration_events_source_id
    ON integration_events (source_id);

CREATE INDEX idx_integration_events_sync_run_id
    ON integration_events (sync_run_id);

CREATE INDEX idx_integration_events_status
    ON integration_events (status);

CREATE INDEX idx_integration_events_received_at
    ON integration_events (received_at DESC);