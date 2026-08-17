-- name: UserHasPermission :one
SELECT EXISTS (
    SELECT 1
    FROM user_groups ug
    JOIN group_permissions gp ON gp.group_id = ug.group_id
    JOIN permissions p ON p.id = gp.permission_id
    WHERE ug.user_id = $1 AND p.resource = $2 AND p.action = $3
) AS has_permission;

-- name: GetGroupByName :one
SELECT * FROM groups WHERE name = $1;

-- name: AddUserToGroup :exec
INSERT INTO user_groups (user_id, group_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;