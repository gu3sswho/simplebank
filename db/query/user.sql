-- name: CreateUser :one
INSERT INTO users (
  username, hashed_password, full_name, email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- -- name: GetUserForUpdate :one
-- SELECT * FROM users
-- WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- -- name: ListUsers :many
-- SELECT * FROM users
-- ORDER BY id
-- LIMIT $1
-- OFFSET $2;

-- -- name: UpdateUser :one
-- UPDATE users
--   set balance = $2
-- WHERE id = $1
-- RETURNING *;

-- -- name: UpdateUserPassword :one
-- UPDATE users
--   set hashed_password = sqlc.arg(hashed_password)
-- WHERE id = sqlc.arg(id)
-- RETURNING *;

-- -- name: DeleteUser :exec
-- DELETE FROM users
-- WHERE id = $1;