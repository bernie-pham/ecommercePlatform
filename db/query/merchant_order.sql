-- name: AddMerchantOrder :one
INSERT INTO merchant_order (
    merchant_id,
    total_price,
    order_status,
    order_id
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateMerchantOrderStatus :exec
UPDATE merchant_order
SET order_status = $1
WHERE id = $2; 