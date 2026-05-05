DROP TABLE IF EXISTS secrets;

CREATE TYPE secret_type AS ENUM ('password', 'bank_card', 'text', 'binary');

CREATE TABLE IF NOT EXISTS secrets (
    id serial primary key,
    user_id INTEGER,
    secret_type secret_type ,
    secret BYTEA,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMP  WITH TIME ZONE,
    comment varchar(200)
);

ALTER TABLE secrets
ADD CONSTRAINT fk_users_secrets
FOREIGN KEY (user_id)
REFERENCES users (user_id)
ON DELETE CASCADE
ON UPDATE RESTRICT;