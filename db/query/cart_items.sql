-- name: AddCartItem :one
INSERT INTO cart_item (
    product_entry_id,
    quantity,
    user_id
)
VALUES ($1, $2, $3) RETURNING *;

-- name: AddDupCartItem :one
UPDATE cart_item
SET
    quantity = quantity + $1,
    modified_at = $2
WHERE 
    id = $3 AND
    user_id = $4
RETURNING *;

-- name: UpdateCartItem :exec
UPDATE cart_item
SET
    quantity = $1,
    modified_at = now()
WHERE 
    id = $2 AND
    user_id = $3;

-- name: GetPEQtyByID :one
SELECT pde.quantity
FROM cart_item c
LEFT JOIN product_entry pde ON c.product_entry_id = pde.id 
WHERE c.id = $1;

-- name: DeleteCartItemByID :exec
DELETE FROM cart_item
WHERE id = $1 AND user_id = $2;

-- name: DeleteAllCartItemByUserID :exec
DELETE FROM cart_item
WHERE user_id = $1;

-- name: ListCartItemsByUserID :many
SELECT *
FROM cart_item
WHERE user_id = $1;

-- name: GetCartItemByEntryID :one
SELECT id 
FROM cart_item
WHERE user_id = $1 AND product_entry_id = $2;

-- name: GetCartItemByID :one
SELECT * 
FROM cart_item
WHERE id = $1;

-- name: GetMerchantByCartID :one
SELECT p.merchant_id
FROM cart_item c
LEFT JOIN product_entry pde ON c.product_entry_id = pde.id
LEFT JOIN products p ON pde.product_id = p.id 
WHERE c.id = $1;

