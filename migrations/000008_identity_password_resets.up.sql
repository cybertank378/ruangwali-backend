BEGIN;

CREATE TABLE password_reset_tokens (
                                       id UUID PRIMARY KEY,

                                       user_id UUID NOT NULL,

                                       token_hash BYTEA NOT NULL,

                                       expires_at TIMESTAMPTZ NOT NULL,

                                       used_at TIMESTAMPTZ NULL,

                                       revoked_at TIMESTAMPTZ NULL,

                                       replaced_by_id UUID NULL,

                                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                       CONSTRAINT fk_password_reset_tokens_user
                                           FOREIGN KEY (user_id)
                                               REFERENCES users(id)
                                               ON DELETE CASCADE,

                                       CONSTRAINT fk_password_reset_tokens_replaced_by
                                           FOREIGN KEY (replaced_by_id)
                                               REFERENCES password_reset_tokens(id)
                                               ON DELETE SET NULL,

                                       CONSTRAINT uq_password_reset_tokens_token_hash
                                           UNIQUE (token_hash),

                                       CONSTRAINT chk_password_reset_tokens_token_hash_not_empty
                                           CHECK (
                                               octet_length(token_hash) > 0
                                               ),

                                       CONSTRAINT chk_password_reset_tokens_expiration
                                           CHECK (
                                               expires_at > created_at
                                               ),

                                       CONSTRAINT chk_password_reset_tokens_used_at
                                           CHECK (
                                               used_at IS NULL
                                                   OR used_at >= created_at
                                               ),

                                       CONSTRAINT chk_password_reset_tokens_revoked_at
                                           CHECK (
                                               revoked_at IS NULL
                                                   OR revoked_at >= created_at
                                               ),

                                       CONSTRAINT chk_password_reset_tokens_updated_at
                                           CHECK (
                                               updated_at >= created_at
                                               ),

                                       CONSTRAINT chk_password_reset_tokens_terminal_state
                                           CHECK (
                                               NOT (
                                                   used_at IS NOT NULL
                                                       AND revoked_at IS NOT NULL
                                                   )
                                               ),

                                       CONSTRAINT chk_password_reset_tokens_replacement
                                           CHECK (
                                               replaced_by_id IS NULL
                                                   OR revoked_at IS NOT NULL
                                               )
);

CREATE UNIQUE INDEX uq_password_reset_tokens_active_user
    ON password_reset_tokens (user_id)
    WHERE
        used_at IS NULL
        AND revoked_at IS NULL;

CREATE INDEX idx_password_reset_tokens_user_id
    ON password_reset_tokens (user_id);

CREATE INDEX idx_password_reset_tokens_expires_at
    ON password_reset_tokens (expires_at);

CREATE INDEX idx_password_reset_tokens_active_expiration
    ON password_reset_tokens (expires_at)
    WHERE
        used_at IS NULL
        AND revoked_at IS NULL;

CREATE INDEX idx_password_reset_tokens_replaced_by_id
    ON password_reset_tokens (replaced_by_id)
    WHERE
        replaced_by_id IS NOT NULL;

COMMIT;