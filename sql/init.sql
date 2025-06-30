CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL
);

-- Добавляем пользователя по умолчанию: email - user, password - user, role - user
INSERT INTO users (email, password_hash, role)
VALUES ('user', '$2a$10$gdBeu07h6hnHAwBrFRrNJeN2zJ8gTRBKUFMkdy5WxdBqnX/dx776K', 'user')
ON CONFLICT (email) DO NOTHING;

