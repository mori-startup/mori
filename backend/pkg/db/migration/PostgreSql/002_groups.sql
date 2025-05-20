-- +migrate Up
CREATE TABLE IF NOT EXISTS groups (
    group_id VARCHAR(100) NOT NULL PRIMARY KEY,
    administrator VARCHAR(50) NOT NULL,
    name VARCHAR(50) NOT NULL,
    description VARCHAR(255)
);

-- +migrate Down
DROP TABLE IF EXISTS groups;