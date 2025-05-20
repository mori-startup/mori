-- +migrate Up
CREATE TABLE IF NOT EXISTS group_users (
    group_id VARCHAR(100) NOT NULL,
    user_id VARCHAR(100) NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS group_users;