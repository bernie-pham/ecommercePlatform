-- name: CreateUser :one
INSERT INTO users (
  full_name,
  email,
  phone,
  hashed_password,
  access_level  
)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
Where email = $1;

-- name: UpdateUser :one
UPDATE users
SET 
  full_name = COALESCE(sqlc.narg(full_name), full_name),
  phone = COALESCE(sqlc.narg(phone), phone),
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_updated_at = COALESCE(sqlc.narg(password_updated_at), password_updated_at)
WHERE 
  email = sqlc.arg(email)
RETURNING *;



