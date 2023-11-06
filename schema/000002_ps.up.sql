CREATE TABLE ps_item
(
    id              UUID DEFAULT uuid_generate_v4() PRIMARY KEY UNIQUE,
    title           varchar(255) not null UNIQUE,
    userlogin       varchar(255) not null,
    password_hash   varchar(255) not null,
    description     varchar(255)
);

CREATE TABLE users_item
(
    id              UUID DEFAULT uuid_generate_v4() PRIMARY KEY UNIQUE,
    user_id         UUID references users (id) on delete cascade not null,
    list_id         UUID references ps_item (id) on delete cascade not null
);