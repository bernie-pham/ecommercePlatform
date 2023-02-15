-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id,
    product_entry_id,
    quantity,
    total_price
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateOrderItemV2 :exec
INSERT INTO order_items (
    order_id,
    product_entry_id,
    quantity,
    total_price,
    merchant_order_id
)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateOrderItemQuantity :one
UPDATE order_items
SET 
    quantity = $1
WHERE 
    order_id = $2 AND product_entry_id = $3
RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE order_id = $1 AND product_entry_id = $2;

-- name: ListOrderItemsByMerchantOrderID :many
SELECT product_entry_id, quantity, total_price
FROM order_items
WHERE merchant_order_id = $1;