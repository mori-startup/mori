-- +migrate Up
CREATE TABLE IF NOT EXISTS group_messages (
    message_id VARCHAR(100) NOT NULL,
    receiver_id VARCHAR(100) NOT NULL,
    is_read INT DEFAULT 0,
    PRIMARY KEY (message_id)
);
-- +migrate Down
DROP TABLE IF EXISTS group_messages;
