-- name: CreateSession :one
INSERT INTO 
    sessions (id, email, refresh_token, user_agent, client_ip, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSession :one
SELECT 
    email, refresh_token, expires_at
FROM 
    sessions
WHERE 
    id = $1 AND is_blocked is false;