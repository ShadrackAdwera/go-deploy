-- name: CreateLog :one
INSERT INTO logs (user_id,description) VALUES ($1,$2) RETURNING *;

-- name: ListLogs :one
SELECT * FROM logs ORDER BY id LIMIT $1 OFFSET $2;

-- name: GetLog :one
SELECT * FROM logs WHERE id = $1;

-- name: DeleteLog :exec
DELETE FROM logs
WHERE id = $1;