// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: products.sql

package db

import (
	"context"
	"database/sql"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (
    name,
    merchant_id,
    status
)
VALUES ($1, $2, $3)
RETURNING id, name, merchant_id, status, created_at
`

type CreateProductParams struct {
	Name       string            `json:"name"`
	MerchantID int32             `json:"merchant_id"`
	Status     NullProductStatus `json:"status"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, createProduct, arg.Name, arg.MerchantID, arg.Status)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.MerchantID,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const getMerchantIDbyPrID = `-- name: GetMerchantIDbyPrID :one
SELECT merchant_id
FROM products
WHERE id = $1
`

func (q *Queries) GetMerchantIDbyPrID(ctx context.Context, id int64) (int32, error) {
	row := q.db.QueryRowContext(ctx, getMerchantIDbyPrID, id)
	var merchant_id int32
	err := row.Scan(&merchant_id)
	return merchant_id, err
}

const listAllProducts = `-- name: ListAllProducts :many
SELECT id, name, merchant_id, status, created_at
FROM products
`

func (q *Queries) ListAllProducts(ctx context.Context) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, listAllProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.MerchantID,
			&i.Status,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProductByMerchantID = `-- name: ListProductByMerchantID :many
SELECT id, name, merchant_id, status, created_at
FROM products
WHERE merchant_id = $1
`

func (q *Queries) ListProductByMerchantID(ctx context.Context, merchantID int32) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, listProductByMerchantID, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.MerchantID,
			&i.Status,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listProductTags = `-- name: ListProductTags :many
SELECT id, name
FROM product_tags
WHERE id in (
    SELECT product_tags_id
    FROM product_tags_products
    WHERE products_id = $1
)
`

func (q *Queries) ListProductTags(ctx context.Context, productsID int64) ([]ProductTag, error) {
	rows, err := q.db.QueryContext(ctx, listProductTags, productsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ProductTag{}
	for rows.Next() {
		var i ProductTag
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProduct = `-- name: UpdateProduct :one
UPDATE products 
SET 
    name = COALESCE($1, name),
    status = $2
WHERE 
    id = $3
RETURNING id, name, merchant_id, status, created_at
`

type UpdateProductParams struct {
	Name   sql.NullString    `json:"name"`
	Status NullProductStatus `json:"status"`
	ID     int64             `json:"id"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, updateProduct, arg.Name, arg.Status, arg.ID)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.MerchantID,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}
