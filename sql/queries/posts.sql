-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT * FROM posts AS P
INNER JOIN feed_follows AS F
    ON P.feed_id = F.feed_id
    AND F.user_id = $1
ORDER BY COALESCE(published_at, NOW()) DESC
LIMIT $2;