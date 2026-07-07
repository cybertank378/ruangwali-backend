CREATE TABLE refresh_sessions (
                                  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                  user_id UUID NOT NULL,

                                  family_id UUID NOT NULL,

                                  token_hash BYTEA NOT NULL,

                                  user_agent TEXT,
                                  ip_address INET,

                                  expires_at TIMESTAMPTZ NOT NULL,

                                  last_used_at TIMESTAMPTZ,

                                  revoked_at TIMESTAMPTZ,
                                  revoke_reason VARCHAR(100),

                                  replaced_by_id UUID,

                                  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                  CONSTRAINT fk_refresh_sessions_user
                                      FOREIGN KEY (user_id)
                                          REFERENCES users(id)
                                          ON DELETE CASCADE,

                                  CONSTRAINT fk_refresh_sessions_replaced_by
                                      FOREIGN KEY (replaced_by_id)
                                          REFERENCES refresh_sessions(id)
                                          ON DELETE SET NULL,

                                  CONSTRAINT uq_refresh_sessions_token_hash
                                      UNIQUE (token_hash),

                                  CONSTRAINT chk_refresh_sessions_expiration
                                      CHECK (
                                          expires_at > created_at
                                          ),

                                  CONSTRAINT chk_refresh_sessions_last_used_at
                                      CHECK (
                                          last_used_at IS NULL
                                              OR last_used_at >= created_at
                                          ),

                                  CONSTRAINT chk_refresh_sessions_revoked_at
                                      CHECK (
                                          revoked_at IS NULL
                                              OR revoked_at >= created_at
                                          ),

                                  CONSTRAINT chk_refresh_sessions_replacement
                                      CHECK (
                                          replaced_by_id IS NULL
                                              OR replaced_by_id <> id
                                          )
);

CREATE INDEX idx_refresh_sessions_user_id
    ON refresh_sessions (user_id);

CREATE INDEX idx_refresh_sessions_family_id
    ON refresh_sessions (family_id);

CREATE INDEX idx_refresh_sessions_user_active
    ON refresh_sessions (
                         user_id,
                         expires_at
        )
    WHERE revoked_at IS NULL;

CREATE INDEX idx_refresh_sessions_family_active
    ON refresh_sessions (
                         family_id,
                         expires_at
        )
    WHERE revoked_at IS NULL;

CREATE INDEX idx_refresh_sessions_expires_at
    ON refresh_sessions (expires_at);