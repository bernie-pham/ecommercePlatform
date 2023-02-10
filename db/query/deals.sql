-- name: GetDealByID :one
SELECT *
FROM deals
WHERE id = $1 AND 
    start_date <= now() AND now() <= end_date;

-- name: ListDealsByMerchantID :many
SELECT * 
FROM deals
WHERE merchant_id = $1;

-- name: CreateDeal :one
INSERT INTO deals (
    name,
    start_date,
    end_date,
    type,
    discount_rate,
    deal_limit
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;