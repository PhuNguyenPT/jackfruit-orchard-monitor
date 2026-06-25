-- name: InsertAirTempHumidReading :exec
INSERT INTO air_temp_humid_readings (addr, temperature, humidity, created_at)
VALUES ($1, $2, $3, $4);
-- name: GetLatestAirTempHumidReadings :many
SELECT DISTINCT ON (addr)
    addr, temperature, humidity, created_at
FROM air_temp_humid_readings
ORDER BY addr, created_at DESC;

-- name: GetAirTempHumidReadingsByAddr :many
SELECT addr, temperature, humidity, created_at
FROM air_temp_humid_readings
WHERE addr = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: DeleteOldAirTempHumidReadings :exec
DELETE FROM air_temp_humid_readings WHERE created_at < $1;

-- name: GetAirTempHumidReadingsByAddrSince :many
SELECT addr, temperature, humidity, created_at
FROM air_temp_humid_readings
WHERE addr = $1 AND created_at > $2
ORDER BY created_at ASC;