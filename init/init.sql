CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    balance INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS history (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    amount INTEGER
);

CREATE INDEX sender
ON history(sender_id);

CREATE INDEX receiver
ON history(receiver_id);