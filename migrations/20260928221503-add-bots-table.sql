CREATE TABLE bots(
    id SERIAL PRIMARY KEY,
    owner_id INT,
    bot_name VARCHAR(128),
    telegram_token VARCHAR(256),
    created_at TIMESTAMPTZ,
);
-- +migrate Up

DROP TABLE bots;
-- +migrate Down
