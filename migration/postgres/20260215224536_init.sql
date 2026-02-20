-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    value VARCHAR(64) NOT NULL UNIQUE,
    description TEXT
);

INSERT INTO roles (id, value)
VALUES (gen_random_uuid(), 'super_admin');

INSERT INTO roles (id, value)
VALUES (gen_random_uuid(), 'admin');

INSERT INTO roles (id, value)
VALUES (gen_random_uuid(), 'user');

CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login VARCHAR(255) NOT NULL UNIQUE,
    email CITEXT NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,

    version INT NOT NULL DEFAULT 1 CHECK (version > 0),

    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,

    deleted_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS accounts_outbox (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(128) NOT NULL DEFAULT 'account',
    event_type VARCHAR(256) NOT NULL,

    payload JSONB NOT NULL,
    headers JSONB DEFAULT '{}'::jsonb,

    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    retry_count INT NOT NULL DEFAULT 0,
    max_retries INT NOT NULL DEFAULT 3,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    -- Индексы для эффективной выборки
    CONSTRAINT accounts_outbox_status_check CHECK (status IN ('pending', 'processing', 'sent', 'failed')),
    CONSTRAINT accounts_outbox_retry_count_check CHECK (retry_count >= 0 AND retry_count <= max_retries)
);

CREATE INDEX IF NOT EXISTS idx_accounts_outbox_status_id
    ON accounts_outbox(status, id);

CREATE INDEX IF NOT EXISTS idx_accounts_outbox_aggregate_id
    ON accounts_outbox(aggregate_id);

CREATE INDEX IF NOT EXISTS idx_accounts_outbox_pending
    ON accounts_outbox(id)
    WHERE status = 'pending';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accounts_outbox;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
