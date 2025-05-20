-- +migrate Up
CREATE TABLE IF NOT EXISTS almost_private (
    user_id VARCHAR(100) NOT NULL,
    post_id VARCHAR(100) NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS almost_private;
