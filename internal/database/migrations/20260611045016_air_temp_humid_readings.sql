-- +goose Up
-- +goose StatementBegin
CREATE TABLE air_temp_humid_readings (
    id          BIGSERIAL    PRIMARY KEY,
    addr        SMALLINT     NOT NULL,
    temperature SMALLINT     NOT NULL,
    humidity    SMALLINT     NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS air_temp_humid_readings;
-- +goose StatementEnd