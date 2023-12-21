CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE credentials
(
    id              UUID DEFAULT uuid_generate_v4() PRIMARY KEY UNIQUE,
    title           varchar(255) not null UNIQUE,
    userlogin       varchar(255) not null,
    password_hash   varchar(255) not null,
    description     varchar(255)
);

CREATE TABLE users_credentials
(
    id              UUID DEFAULT uuid_generate_v4() PRIMARY KEY UNIQUE,
    user_id         UUID not null,
    list_id         UUID references credentials (id) on delete cascade not null
);