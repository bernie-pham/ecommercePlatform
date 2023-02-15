-- name: CreatePEntry :one
INSERT INTO 
    product_entry (
        product_id,
        colour_id,
        size_id,
        general_criteria_id,
        quantity,
        deal_id
    )
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING *;

-- name: ListPEntriesByPID :many
SELECT *
FROM product_entry
WHERE product_id = $1;

-- name: ListActivePEntriesByPID :many
SELECT *
FROM product_entry
WHERE product_id = $1 AND is_active = true;

-- name: GetPEntry :one
SELECT *
FROM product_entry
WHERE id = $1;

-- name: UpdatePEntry :one
UPDATE product_entry
SET 
    quantity = COALESCE(sqlc.narg(quantity), quantity),
    deal_id = COALESCE(sqlc.narg(deal_id), deal_id),
    modified_at = COALESCE(sqlc.narg(modified_at), modified_at),
    is_active = sqlc.arg(is_active)
WHERE
    id = sqlc.arg(id)
RETURNING *;


-- name: UpdateEntryQuantity :exec
UPDATE product_entry
SET 
    quantity = quantity - $1
WHERE
    id = $2;

-- name: GetMerchantIDByPEntry :one
SELECT p.merchant_id
FROM product_entry pe 
LEFT JOIN products p ON pe.product_id = p.id
WHERE pe.id = $1;