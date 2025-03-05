-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS(
    INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
    VALUES (
        NOW(),
        NOW(),
        $1,
        $2
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follows WHERE user_id = $1;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows WHERE feed_follows.user_id = $1
AND feed_id = (SELECT id from feeds WHERE url = $2);