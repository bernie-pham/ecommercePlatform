-- name: CreateProduct :one
INSERT INTO products (
    name,
    merchant_id,
    status
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListProductByMerchantID :many
SELECT *
FROM products
WHERE merchant_id = $1;

-- name: UpdateProduct :one
UPDATE products 
SET 
    name = COALESCE(sqlc.narg(name), name),
    status = sqlc.arg(status)
WHERE 
    id = sqlc.arg(id)
RETURNING *;

-- name: ListProductTags :many
SELECT *
FROM product_tags
WHERE id in (
    SELECT product_tags_id
    FROM product_tags_products
    WHERE products_id = $1
);

-- name: GetMerchantIDbyPrID :one
SELECT merchant_id
FROM products
WHERE id = $1;


-- name: ListAllProducts :many
SELECT *
FROM products
LIMIT $1
OFFSET $2;


-- name: ListProductID :many
SELECT id 
FROM products
LIMIT $1
OFFSET $2;

-- name: GetProductNameByID :one
SELECT name 
FROM products
WHERE id = $1;