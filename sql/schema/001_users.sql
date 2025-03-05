-- +goose Up

-- This table stores user information
create table users (
    id UUID primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name TEXT UNIQUE NOT NULL
);

-- +goose Down
drop table users;