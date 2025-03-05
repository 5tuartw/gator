-- name: CreatePost :exec
INSERT INTO posts (created_at, updated_at, title, url, description, published_at, feed_id, id)
VALUES (
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: GetUserPosts :many
SELECT posts.*, feeds.name as feed_name
FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
JOIN feeds ON posts.feed_id = feeds.id
JOIN users ON feed_follows.user_id = users.id
WHERE users.name = $1
ORDER BY posts.published_at DESC
LIMIT $2;