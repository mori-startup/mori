-- +migrate Up
CREATE TABLE IF NOT EXISTS followers (
    follower_id VARCHAR(100) NOT NULL,
    user_id VARCHAR(100) NOT NULL,
    CONSTRAINT pk_followers PRIMARY KEY (follower_id, user_id)
);

-- +migrate Down
DROP TABLE IF EXISTS followers;