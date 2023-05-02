-- name: CreateTag :one
INSERT INTO 
    product_tags (name)
VALUES ($1) RETURNING *;

-- name: DeleteTag :exec
DELETE FROM product_tags
WHERE id = $1;

-- name: CreateProTag :one
INSERT INTO 
    product_tags_products (product_tags_id, products_id)
VALUES ($1, $2) RETURNING *;

-- name: DeleteProTag :exec
DELETE FROM product_tags_products
WHERE product_tags_id = $1 AND products_id = $2;

-- name: ListTags :many
SELECT *
FROM product_tags;

-- name: ListProTags :many
SELECT *
FROM product_tags_products;

-- name: ListProductsByTagID :many
SELECT *
FROM products
WHERE id in (
    SELECT products_id
    FROM product_tags_products
    WHERE product_tags_id = $1
);

-- name: ListTagID :many
SELECT id
FROM product_tags
LIMIT $1
OFFSET $2;

-- name: ListProductIDbyTagID :many
SELECT products_id
FROM product_tags_products
WHERE product_tags_id = $1;

-- name: GetTagNameByID :one
SELECT name 
FROM product_tags
WHERE id = $1;