-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(100) NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    email VARCHAR(50) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    nickname VARCHAR(50),
    birthday DATE NOT NULL,
    image VARCHAR(255),
    about TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'PUBLIC',
    password VARCHAR(100) NOT NULL,
    verification_token VARCHAR(100),
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    reset_token VARCHAR(100),
    reset_token_expires TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS users;
