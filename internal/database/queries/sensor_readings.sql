-- name: InsertSensorReading :exec
INSERT INTO sensor_readings (addr, temperature, humidity)
VALUES ($1, $2, $3);

-- name: GetLatestReadings :many
SELECT DISTINCT ON (addr)
    addr, temperature, humidity, created_at
FROM sensor_readings
ORDER BY addr, created_at DESC;

-- name: GetReadingsByAddr :many
SELECT addr, temperature, humidity, created_at
FROM sensor_readings
WHERE addr = $1
ORDER BY created_at DESC
LIMIT $2;