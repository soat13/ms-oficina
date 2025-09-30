-- +migrate Up

CREATE TABLE users
(
    id            UUID PRIMARY KEY,
    name          VARCHAR(255)       NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255)       NOT NULL,
    created_at    TIMESTAMP          NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP          NOT NULL DEFAULT NOW()
);

insert into users (id, name, email, password_hash) values
(gen_random_uuid(), 'Admin', 'admin@admin.com', '$2a$10$7vp9huEQa8T08CSb9Wxpx./iY6KQeoIUMTsZC/344QJk9./UytDie'); -- password: secret123

-- +migrate Down
DROP TABLE IF EXISTS users;
