-- +migrate Up
CREATE TABLE IF NOT EXISTS notifications (
    notif_id VARCHAR(100) NOT NULL,
    user_id VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    content VARCHAR(255) NOT NULL,
    sender VARCHAR(50) NOT NULL,
    PRIMARY KEY (notif_id)
);


-- +migrate Down
DROP TABLE IF EXISTS notifications;