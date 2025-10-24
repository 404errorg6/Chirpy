-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
RETURNING *;

-- name: ResetDatabase :exec
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :exec
UPDATE users
SET updated_at = NOW(), email = $2, hashed_password = $3
WHERE id = $1;
