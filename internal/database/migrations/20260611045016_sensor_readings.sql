-- +goose Up
-- +goose StatementBegin
CREATE TABLE sensor_readings (
    id          BIGSERIAL    PRIMARY KEY,
    addr        TEXT         NOT NULL,
    temperature REAL         NOT NULL,
    humidity    REAL         NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sensor_readings;
-- +goose StatementEnd