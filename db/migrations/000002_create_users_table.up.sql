CREATE TABLE IF NOT EXISTS users (
    id          VARCHAR(27)        PRIMARY KEY,
    username    VARCHAR(255)       NOT NULL,
    email       VARCHAR(255)       UNIQUE NOT NULL,
    password    VARCHAR(255)       NOT NULL,
    created_at  BIGINT             NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()),
    updated_at  BIGINT,
    deleted_at  BIGINT NULL
);
