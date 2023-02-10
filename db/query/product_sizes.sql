-- name: CreatePSize :one
INSERT INTO 
    product_size (size_value)
VALUES ($1) RETURNING *;

-- name: DeletePSize :exec
DELETE FROM product_size
WHERE id = $1;

-- name: ListPSizes :many
SELECT *
FROM product_size;