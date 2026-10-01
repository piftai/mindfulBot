-- +migrate Up
CREATE TABLE IF NOT EXISTS bots(
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    st_id INT NOT NULL UNIQUE,
    owner_id INT NOT NULL,
    bot_name VARCHAR(128) NOT NULL,
    telegram_token VARCHAR(256) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_bots_st_id ON bots (st_id);

-- +migrate Down
DROP TABLE IF EXISTS bots;
