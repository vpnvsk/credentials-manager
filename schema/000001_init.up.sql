CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO users (id, username, password_hash)
VALUES (uuid_generate_v4(), 'new_user', 'password_hash');

CREATE TABLE users
(
    id            UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username      varchar(255) not null unique,
    password_hash varchar(255) not null
);