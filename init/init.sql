CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    balance INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS history (
    id SERIAL PRIMARY KEY,
    sender_name TEXT NOT NULL,
    receiver_name TEXT NOT NULL,
    amount INTEGER
);

CREATE TABLE IF NOT EXISTS inventory (
    id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    item TEXT NOT NULL
);

CREATE INDEX sender
    ON history(sender_name);

CREATE INDEX receiver
    ON history(receiver_name);

CREATE INDEX idx_owner_id_item
    ON inventory (owner_id, item);