INSERT INTO users (user_id, username, password_hash, created_at)
VALUES (0, 'username', 'hash', '2024-12-25 10:30:00'),
       (10, 'username_10', 'hash', '2024-12-25 10:30:00')
       ;


INSERT INTO secrets (id, user_id, secret_type, secret, created_at, updated_at, comment)
VALUES (
        0, 0, 'text',
        '\xDEADBEEF'::bytea, '2024-12-25 10:30:00',
        '2024-12-25 10:30:00', 'comment'
);