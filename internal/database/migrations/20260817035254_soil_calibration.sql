-- +goose Up
-- +goose StatementBegin
CREATE TABLE soil_calibration (
    sensor_idx SMALLINT     PRIMARY KEY,
    dry_value  SMALLINT     NOT NULL,
    wet_value  SMALLINT     NOT NULL,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS soil_calibration;
-- +goose StatementEnd
