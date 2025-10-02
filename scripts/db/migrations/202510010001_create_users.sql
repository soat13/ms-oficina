-- +migrate Up
CREATE TYPE user_role AS ENUM ('attendant', 'manager', 'mechanic');

CREATE TABLE users
(
    id            UUID         PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    document      VARCHAR(20)  UNIQUE NOT NULL,
    document_type VARCHAR(4)   NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    phone_number  VARCHAR(11)  NOT NULL,
    password      VARCHAR(255) NOT NULL,
    roles         user_role[]  NOT NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_roles ON users USING GIN (roles);

-- +migrate Down
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_role;

