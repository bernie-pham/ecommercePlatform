-- name: CreateOrder :one
INSERT INTO orders (user_id, base_price, discount_price)
VALUES ($1, $2, $3) RETURNING id;

-- name: UpdateOrder :one
UPDATE orders 
SET 
    status = COALESCE(sqlc.narg(status):: order_status, status)  ,
    deal_id = COALESCE(sqlc.narg(deal_id), deal_id),
    base_price = COALESCE(sqlc.narg(base_price), base_price),
    discount_price = COALESCE(sqlc.narg(discount_price), discount_price)
WHERE
    id = sqlc.arg(id)
RETURNING *;

-- name: ListOrder :many
SELECT * 
FROM orders 
WHERE user_id = $1;

-- name: GetCurrentOrder :one
SELECT *
FROM orders
WHERE user_id = $1 AND status = 'open';