-- name: CreateVerification :one
INSERT INTO
    verifications (
        id,
        email,
        expires_at
    )
VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetVerification :one
SELECT *
FROM verifications
WHERE id = $1 AND is_occurpied = false;

-- name: OccupyVerification :exec
UPDATE verifications
SET is_occurpied = true
Where id = $1 AND is_occurpied = false;
