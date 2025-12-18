DROP TABLE IF EXISTS secrets;

CREATE TABLE IF NOT EXISTS secrets (
    id serial primary key,
    user_id INTEGER,
    secret_type varchar(10) not null ,
    secret BYTEA,
    created_at TIMESTAMP,
    comment varchar(200)
);

ALTER TABLE secrets
ADD CONSTRAINT fk_users_secrets
FOREIGN KEY (user_id)
REFERENCES users (user_id)
ON DELETE CASCADE
ON UPDATE RESTRICT;