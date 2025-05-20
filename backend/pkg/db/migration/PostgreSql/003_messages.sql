-- +migrate Up
CREATE TABLE IF NOT EXISTS messages (
    message_id VARCHAR(100) NOT NULL PRIMARY KEY,
    sender_id VARCHAR(100) NOT NULL,
    receiver_id VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    content VARCHAR(255) NOT NULL,
    is_read INT DEFAULT 0
);

-- +migrate Down
DROP TABLE IF EXISTS messages;