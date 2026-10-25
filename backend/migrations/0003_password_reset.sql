CREATE TABLE password_reset_tokens (
    id         uuid PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash bytea NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    used_at    timestamptz,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL
);
CREATE INDEX password_reset_tokens_user_idx ON password_reset_tokens(user_id);
GRANT SELECT, INSERT, UPDATE, DELETE ON password_reset_tokens TO ssp_app;
