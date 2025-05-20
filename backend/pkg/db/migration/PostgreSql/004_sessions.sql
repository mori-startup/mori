-- +migrate Up
CREATE TABLE IF NOT EXISTS sessions (
    session_id VARCHAR(100) NOT NULL PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,
    expiration_time TIMESTAMP NOT NULL
);
-- +migrate Down
DROP TABLE "sessions";