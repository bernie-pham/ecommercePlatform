-- name: CreatePPrice :one
INSERT INTO 
    product_pricing (
        product_id,
        base_price,
        start_date,
        end_date,
        priority
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: GetTodayBasePrice :one
SELECT base_price
FROM product_pricing
WHERE (product_id, priority) in (
    select pp.product_id, MAX(pp.priority)
	from product_pricing pp
	where pp.product_id = $1 and pp.start_date <= now() AND now() <= pp.end_date AND 
	    pp.is_active = true
	group by pp.product_id
);

-- name: UpdatePPrice :one
UPDATE product_pricing
SET 
    base_price = COALESCE(sqlc.narg(base_price), base_price),
    end_date = COALESCE(sqlc.narg(end_date), end_date),
    is_active = sqlc.arg(is_active),
    priority = COALESCE(sqlc.narg(priority), priority)
WHERE 
    id = sqlc.arg(id)
RETURNING *;

-- name: ListPriceByPID :many
SELECT *
FROM product_pricing
WHERE product_id = $1;