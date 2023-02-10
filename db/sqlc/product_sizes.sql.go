// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: product_sizes.sql

package db

import (
	"context"
)

const createPSize = `-- name: CreatePSize :one
INSERT INTO 
    product_size (size_value)
VALUES ($1) RETURNING id, size_value
`

func (q *Queries) CreatePSize(ctx context.Context, sizeValue string) (ProductSize, error) {
	row := q.db.QueryRowContext(ctx, createPSize, sizeValue)
	var i ProductSize
	err := row.Scan(&i.ID, &i.SizeValue)
	return i, err
}

const deletePSize = `-- name: DeletePSize :exec
DELETE FROM product_size
WHERE id = $1
`

func (q *Queries) DeletePSize(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePSize, id)
	return err
}

const listPSizes = `-- name: ListPSizes :many
SELECT id, size_value
FROM product_size
`

func (q *Queries) ListPSizes(ctx context.Context) ([]ProductSize, error) {
	rows, err := q.db.QueryContext(ctx, listPSizes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ProductSize{}
	for rows.Next() {
		var i ProductSize
		if err := rows.Scan(&i.ID, &i.SizeValue); err != nil {
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