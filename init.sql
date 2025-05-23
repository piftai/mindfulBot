CREATE TABLE admin (
                       id SERIAL PRIMARY KEY,
                       user_id BIGINT NOT NULL,
                       username TEXT
);

CREATE TABLE reminders (
                           id SERIAL PRIMARY KEY,
                           user_id BIGINT NOT NULL,
                           username TEXT,
                           day TEXT NOT NULL,
                           time TEXT NOT NULL,
                           remind_1h TIMESTAMP NOT NULL,
                           remind_24h TIMESTAMP NOT NULL,
                           is_always BOOLEAN DEFAULT FALSE
);

CREATE TABLE users (
                       user_id BIGINT PRIMARY KEY,
                       username TEXT,
                       id SERIAL
);