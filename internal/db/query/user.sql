-- name: CreateUser :one
INSERT INTO users (username,email,tenant_id, password) 
VALUES ($1,$2,$3, $4) 
RETURNING *;

-- name: ListUsers :many
SELECT (users.id, users.username, users.email, tenant.id, tenant.name, tenant.logo) 
FROM users
JOIN tenant 
ON tenant.id = users.tenant_id 
ORDER BY users.id
LIMIT $1 
OFFSET $2;

-- name: GetUser :one
SELECT (users.id, users.username, users.email, tenant.id, tenant.name, tenant.logo) 
FROM users
JOIN tenant 
ON tenant.id = users.tenant_id
WHERE users.id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET
  username = COALESCE(sqlc.narg(username),username),
  email = COALESCE(sqlc.narg(email),email),
  tenant_id = COALESCE(sqlc.narg(tenant_id),tenant_id),
  password = COALESCE(sqlc.narg(password),password)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;