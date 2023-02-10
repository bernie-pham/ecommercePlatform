// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: product_tags.sql

package db

import (
	"context"
)

const createProTag = `-- name: CreateProTag :one
INSERT INTO 
    product_tags_products (product_tags_id, products_id)
VALUES ($1, $2) RETURNING product_tags_id, products_id
`

type CreateProTagParams struct {
	ProductTagsID int64 `json:"product_tags_id"`
	ProductsID    int64 `json:"products_id"`
}

func (q *Queries) CreateProTag(ctx context.Context, arg CreateProTagParams) (ProductTagsProduct, error) {
	row := q.db.QueryRowContext(ctx, createProTag, arg.ProductTagsID, arg.ProductsID)
	var i ProductTagsProduct
	err := row.Scan(&i.ProductTagsID, &i.ProductsID)
	return i, err
}

const createTag = `-- name: CreateTag :one
INSERT INTO 
    product_tags (name)
VALUES ($1) RETURNING id, name
`

func (q *Queries) CreateTag(ctx context.Context, name string) (ProductTag, error) {
	row := q.db.QueryRowContext(ctx, createTag, name)
	var i ProductTag
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteProTag = `-- name: DeleteProTag :exec
DELETE FROM product_tags_products
WHERE product_tags_id = $1 AND products_id = $2
`

type DeleteProTagParams struct {
	ProductTagsID int64 `json:"product_tags_id"`
	ProductsID    int64 `json:"products_id"`
}

func (q *Queries) DeleteProTag(ctx context.Context, arg DeleteProTagParams) error {
	_, err := q.db.ExecContext(ctx, deleteProTag, arg.ProductTagsID, arg.ProductsID)
	return err
}

const deleteTag = `-- name: DeleteTag :exec
DELETE FROM product_tags
WHERE id = $1
`

func (q *Queries) DeleteTag(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTag, id)
	return err
}

const listProTags = `-- name: ListProTags :many
SELECT product_tags_id, products_id
FROM product_tags_products
`

func (q *Queries) ListProTags(ctx context.Context) ([]ProductTagsProduct, error) {
	rows, err := q.db.QueryContext(ctx, listProTags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ProductTagsProduct{}
	for rows.Next() {
		var i ProductTagsProduct
		if err := rows.Scan(&i.ProductTagsID, &i.ProductsID); err != nil {
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

const listProductsByTagID = `-- name: ListProductsByTagID :many
SELECT id, name, merchant_id, status, created_at
FROM products
WHERE id in (
    SELECT products_id
    FROM product_tags_products
    WHERE product_tags_id = $1
)
`

func (q *Queries) ListProductsByTagID(ctx context.Context, productTagsID int64) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, listProductsByTagID, productTagsID)
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

const listTags = `-- name: ListTags :many
SELECT id, name
FROM product_tags
`

func (q *Queries) ListTags(ctx context.Context) ([]ProductTag, error) {
	rows, err := q.db.QueryContext(ctx, listTags)
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