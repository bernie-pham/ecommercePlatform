-- name: CreateMerchantOrder :one
INSERT INTO merchant_order (
    merchant_id,
    total_price,
    order_status,
    order_id
)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: UpdateMerchantOrderStatus :one
UPDATE merchant_order
SET order_status = $1
WHERE id = $2 AND merchant_id = $3
RETURNING *;

-- name: UpdateMerchantOrderTotalPrice :exec
UPDATE merchant_order
SET total_price = $1
WHERE id = $2;

-- name: GetMerchantOrder :one
SELECT *
FROM merchant_order
WHERE id = $1 AND merchant_id = $2;
