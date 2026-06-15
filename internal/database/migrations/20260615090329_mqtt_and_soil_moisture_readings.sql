-- +goose Up
-- +goose StatementBegin

CREATE TABLE mqtt_credentials (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    username    TEXT        NOT NULL UNIQUE,
    password    TEXT        NOT NULL,       -- bcrypt hash
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mqtt_acl (
    id            UUID        PRIMARY KEY DEFAULT uuidv7(),
    credential_id UUID        NOT NULL REFERENCES mqtt_credentials(id) ON DELETE CASCADE,
    topic         TEXT        NOT NULL,     -- e.g. 'sht40/+/data'
    permission    TEXT        NOT NULL DEFAULT 'r' CHECK (permission IN ('r', 'w', 'rw')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (credential_id, topic)
);

CREATE TABLE soil_moisture_readings (
    id          BIGSERIAL   PRIMARY KEY,
    sensor_idx  SMALLINT    NOT NULL,
    raw         SMALLINT    NOT NULL,       -- ADC raw 0–4095
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS soil_moisture_readings;
DROP TABLE IF EXISTS mqtt_acl;
DROP TABLE IF EXISTS mqtt_credentials;
-- +goose StatementEnd