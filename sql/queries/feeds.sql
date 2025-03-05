-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedId :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeedName :one
SELECT name FROM feeds WHERE id = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetched_at NULLS FIRST, last_fetched_at LIMIT 1;