-- +goose Up

-- This table stores the posts
create table posts (
    id UUID primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    title text,
    url text unique not null,
    description text,
    published_at timestamp not null,
    feed_id UUID not null references feeds(id)
);

-- +goose Down
drop table posts;