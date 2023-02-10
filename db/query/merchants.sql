-- name: CreateMerchant :one
INSERT INTO merchants (
  user_id,
  country_code,
  merchant_name,
  description
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListMerchants :many
SELECT *
FROM merchants
Where user_id = $1;

-- name: GetMerchant :one
SELECT *
FROM merchants
WHERE id = $1;

-- name: DisableMerchant :exec
UPDATE merchants
SET 
    is_active = false
WHERE id = $1;

-- name: EnableMerchant :exec
UPDATE merchants
SET 
    is_active = true
WHERE id = $1;

-- name: UpdateMerchant :one
UPDATE merchants
SET 
  merchant_name = COALESCE(sqlc.narg(merchant_name), merchant_name),
  country_code = COALESCE(sqlc.narg(country_code), country_code),
  description = COALESCE(sqlc.narg(description), description)
WHERE 
  id = sqlc.arg(id)
RETURNING *;



