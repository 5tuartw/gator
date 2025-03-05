-- +goose Up

-- This is a joining table that records who is following different feeds
create table feed_follows (
    id serial primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id UUID references users(id) on delete cascade not null,
    feed_id UUID references feeds(id) on delete cascade not null,
    unique (user_id, feed_id) 
);

-- +goose Down
drop table feed_follows;