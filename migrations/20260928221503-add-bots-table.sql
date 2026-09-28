-- +migrate Up
CREATE TABLE IF NOT EXISTS bots(
    id SERIAL PRIMARY KEY,
    st_id INT,
    owner_id INT,
    bot_name VARCHAR(128),
    telegram_token VARCHAR(256),
    created_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_bots_st_id ON bots (st_id);

-- +migrate Down
DROP INDEX IF EXISTS uq_bots_st_id;
DROP TABLE IF EXISTS bots;
