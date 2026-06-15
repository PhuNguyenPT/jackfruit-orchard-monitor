-- name: GetMQTTCredentialByUsername :one
SELECT * FROM mqtt_credentials WHERE username = $1;

-- name: CreateMQTTCredential :one
INSERT INTO mqtt_credentials (username, password)
VALUES ($1, $2)
RETURNING *;

-- name: GetMQTTACLByCredentialID :many
SELECT * FROM mqtt_acl WHERE credential_id = $1;

-- name: CreateMQTTACL :one
INSERT INTO mqtt_acl (credential_id, topic, permission)
VALUES ($1, $2, $3)
RETURNING *;