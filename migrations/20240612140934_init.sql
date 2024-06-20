-- +goose Up
CREATE TABLE IF NOT EXISTS configurations
(
    id serial primary key,
    name TEXT NOT NULL,
    code varchar(4) NOT NULL,
    is_available boolean NOT NULL,
    type smallint NOT NULL,
    created_at TIMESTAMPTZ DEFAULT current_timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS pairs
(
    id serial primary key,
    "from" varchar(4) NOT NULL,
    "to" varchar(4) NOT NULL,
    amount decimal NOT NULL,
    created_at TIMESTAMPTZ DEFAULT current_timestamp NOT NULL
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS configurations_code_idx ON configurations (code);
CREATE INDEX CONCURRENTLY IF NOT EXISTS pairs_from_to_idx ON pairs ("from", "to");

-- +goose Down

