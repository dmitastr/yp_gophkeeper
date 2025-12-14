CREATE TABLE IF NOT EXISTS users (
    user_id serial primary key,
    username VARCHAR(40) not null,
    password_hash VARCHAR(40) not null,
    created TIMESTAMP
);