-- +goose Up
-- +goose StatementBegin
CREATE TABLE groups (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    name        TEXT        NOT NULL UNIQUE,   -- e.g. 'admin', 'user'
    description TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE permissions (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    resource    TEXT        NOT NULL,          -- e.g. 'soil_calibration', 'sensors', 'users'
    action      CHAR(1)     NOT NULL CHECK (action IN ('r', 'w', 'd')),
    description TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (resource, action)
);

CREATE TABLE group_permissions (
    group_id      UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, permission_id)
);

CREATE TABLE user_groups (
    user_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, group_id)
);

-- Seed the two starting groups.
INSERT INTO groups (name, description) VALUES
    ('admin', 'Full access'),
    ('user',  'Limited access');

-- Seed the initial permission set: resource + r/w/d action.
INSERT INTO permissions (resource, action, description) VALUES
    ('soil_calibration', 'r', 'View per-sensor soil calibration values'),
    ('soil_calibration', 'w', 'Edit per-sensor soil calibration values'),
    ('sensors',          'r', 'View sensor dashboards and history');

-- admin gets every permission that exists at migration time.
INSERT INTO group_permissions (group_id, permission_id)
SELECT g.id, p.id FROM groups g CROSS JOIN permissions p WHERE g.name = 'admin';

-- user gets read-only sensor access, no calibration write.
INSERT INTO group_permissions (group_id, permission_id)
SELECT g.id, p.id FROM groups g JOIN permissions p
    ON p.resource = 'sensors' AND p.action = 'r'
WHERE g.name = 'user';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_groups;
DROP TABLE IF EXISTS group_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS groups;
-- +goose StatementEnd