CREATE TABLE admin (
                id SERIAL PRIMARY KEY,
                user_id BIGINT NOT NULL,
                username TEXT
);