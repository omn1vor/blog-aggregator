-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds 
ORDER BY name;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY COALESCE(last_fetched_at, '0001-01-01')
LIMIT $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET 
    last_fetched_at = NOW(),
    updated_at = NOW()
WHERE id = $1;