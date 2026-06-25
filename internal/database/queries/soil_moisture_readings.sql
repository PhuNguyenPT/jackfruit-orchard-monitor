-- name: InsertSoilMoistureReading :exec
INSERT INTO soil_moisture_readings (sensor_idx, raw, created_at)
VALUES ($1, $2, $3);

-- name: GetLatestSoilMoistureReadings :many
SELECT DISTINCT ON (sensor_idx) sensor_idx, raw, created_at
FROM soil_moisture_readings
ORDER BY sensor_idx, created_at DESC;

-- name: DeleteOldSoilMoistureReadings :exec
DELETE FROM soil_moisture_readings WHERE created_at < $1;

-- name: GetSoilMoistureReadingsBySensorIdx :many
SELECT sensor_idx, raw, created_at
FROM soil_moisture_readings
WHERE sensor_idx = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: GetSoilMoistureReadingsBySensorIdxSince :many
SELECT sensor_idx, raw, created_at
FROM soil_moisture_readings
WHERE sensor_idx = $1 AND created_at > $2
ORDER BY created_at ASC;