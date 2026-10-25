-- SSP-52: initial schema. Enum values are enforced with CHECK constraints.
CREATE TABLE users (
    id                    uuid PRIMARY KEY,
    email                 text NOT NULL UNIQUE CHECK (email = lower(email)),
    password_hash         text NOT NULL,
    first_name            text NOT NULL,
    last_name             text NOT NULL,
    role                  text NOT NULL CHECK (role IN ('parent','specialist','admin')),
    status                text NOT NULL CHECK (status IN ('active','inactive')),
    specialist_approved   boolean NOT NULL DEFAULT false,
    failed_login_attempts integer NOT NULL DEFAULT 0,
    locked_until          timestamptz,
    created_at            timestamptz NOT NULL,
    updated_at            timestamptz NOT NULL
);

CREATE TABLE refresh_tokens (
    id         uuid PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    family_id  uuid NOT NULL,
    token_hash bytea NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    rotated_at timestamptz,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL
);
CREATE INDEX refresh_tokens_family_idx ON refresh_tokens(family_id);
CREATE INDEX refresh_tokens_user_idx ON refresh_tokens(user_id);

CREATE TABLE children (
    id         uuid PRIMARY KEY,
    first_name text NOT NULL,
    last_name  text NOT NULL,
    birth_date date NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);

CREATE TABLE child_guardians (
    child_id     uuid NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    guardian_id  uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    relationship text NOT NULL CHECK (relationship IN ('parent','guardian')),
    is_primary   boolean NOT NULL,
    created_at   timestamptz NOT NULL,
    PRIMARY KEY (child_id, guardian_id)
);
CREATE INDEX child_guardians_guardian_idx ON child_guardians(guardian_id);
CREATE UNIQUE INDEX child_guardians_one_primary ON child_guardians(child_id) WHERE is_primary;

CREATE TABLE child_specialists (
    child_id      uuid NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    specialist_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_by   uuid NOT NULL REFERENCES users(id),
    assigned_at   timestamptz NOT NULL,
    ended_at      timestamptz,
    PRIMARY KEY (child_id, specialist_id)
);
CREATE INDEX child_specialists_specialist_idx ON child_specialists(specialist_id);

CREATE TABLE consent_records (
    id             uuid PRIMARY KEY,
    child_id       uuid NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    guardian_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    policy_version text NOT NULL,
    granted_at     timestamptz NOT NULL
);

CREATE TABLE support_sessions (
    id            uuid PRIMARY KEY,
    child_id      uuid NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    specialist_id uuid NOT NULL REFERENCES users(id),
    created_by    uuid NOT NULL REFERENCES users(id),
    title         text NOT NULL,
    starts_at     timestamptz NOT NULL,
    ends_at       timestamptz NOT NULL CHECK (ends_at > starts_at),
    mode          text NOT NULL CHECK (mode IN ('online','in_person')),
    status        text NOT NULL CHECK (status IN ('scheduled','completed','cancelled')),
    created_at    timestamptz NOT NULL,
    updated_at    timestamptz NOT NULL
);
CREATE INDEX support_sessions_child_idx ON support_sessions(child_id, starts_at);

CREATE TABLE diary_entries (
    id          uuid PRIMARY KEY,
    child_id    uuid NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    author_id   uuid NOT NULL REFERENCES users(id),
    type        text NOT NULL CHECK (type IN ('observation','progress','recommendation')),
    visibility  text NOT NULL CHECK (visibility IN ('shared','specialists_only')),
    content     text NOT NULL,
    observed_at timestamptz NOT NULL,
    created_at  timestamptz NOT NULL,
    updated_at  timestamptz NOT NULL
);
CREATE INDEX diary_entries_access_idx ON diary_entries(child_id, visibility, observed_at DESC, id DESC);

CREATE TABLE diary_entry_versions (
    id             uuid PRIMARY KEY,
    entry_id       uuid NOT NULL REFERENCES diary_entries(id) ON DELETE CASCADE,
    version_number integer NOT NULL,
    changed_by     uuid NOT NULL REFERENCES users(id),
    type           text NOT NULL,
    visibility     text NOT NULL,
    content        text NOT NULL,
    observed_at    timestamptz NOT NULL,
    changed_at     timestamptz NOT NULL,
    UNIQUE (entry_id, version_number)
);

-- Audit rows outlive the child on purpose (no FK); retention per ADR-0005.
CREATE TABLE audit_events (
    id          uuid PRIMARY KEY,
    actor_id    uuid,
    child_id    uuid,
    action      text NOT NULL,
    resource    text NOT NULL,
    resource_id uuid,
    outcome     text NOT NULL CHECK (outcome IN ('success','denied','failure')),
    request_id  text NOT NULL,
    occurred_at timestamptz NOT NULL
);
CREATE INDEX audit_events_child_idx ON audit_events(child_id, occurred_at DESC, id DESC);
CREATE INDEX audit_events_time_idx ON audit_events(occurred_at DESC, id DESC);
