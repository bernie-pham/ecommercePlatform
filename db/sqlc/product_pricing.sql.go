// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: product_pricing.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createPPrice = `-- name: CreatePPrice :one
INSERT INTO 
    product_pricing (
        product_id,
        base_price,
        start_date,
        end_date,
        priority
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING id, product_id, base_price, start_date, end_date, is_active, priority
`

type CreatePPriceParams struct {
	ProductID string    `json:"product_id"`
	BasePrice float32   `json:"base_price"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Priority  int32     `json:"priority"`
}

func (q *Queries) CreatePPrice(ctx context.Context, arg CreatePPriceParams) (ProductPricing, error) {
	row := q.db.QueryRowContext(ctx, createPPrice,
		arg.ProductID,
		arg.BasePrice,
		arg.StartDate,
		arg.EndDate,
		arg.Priority,
	)
	var i ProductPricing
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.BasePrice,
		&i.StartDate,
		&i.EndDate,
		&i.IsActive,
		&i.Priority,
	)
	return i, err
}

const getTodayBasePrice = `-- name: GetTodayBasePrice :one
SELECT base_price
FROM product_pricing
WHERE (product_id, priority) in (
    select pp.product_id, MAX(pp.priority)
	from product_pricing pp
	where pp.product_id = $1 and pp.start_date <= now() AND now() <= pp.end_date AND 
	    pp.is_active = true
	group by pp.product_id
)
`

func (q *Queries) GetTodayBasePrice(ctx context.Context, productID string) (float32, error) {
	row := q.db.QueryRowContext(ctx, getTodayBasePrice, productID)
	var base_price float32
	err := row.Scan(&base_price)
	return base_price, err
}

const listPriceByPID = `-- name: ListPriceByPID :many
SELECT id, product_id, base_price, start_date, end_date, is_active, priority
FROM product_pricing
WHERE product_id = $1
`

func (q *Queries) ListPriceByPID(ctx context.Context, productID string) ([]ProductPricing, error) {
	rows, err := q.db.QueryContext(ctx, listPriceByPID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ProductPricing{}
	for rows.Next() {
		var i ProductPricing
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.BasePrice,
			&i.StartDate,
			&i.EndDate,
			&i.IsActive,
			&i.Priority,
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

const updatePPrice = `-- name: UpdatePPrice :one
UPDATE product_pricing
SET 
    base_price = COALESCE($1, base_price),
    end_date = COALESCE($2, end_date),
    is_active = $3,
    priority = COALESCE($4, priority)
WHERE 
    id = $5
RETURNING id, product_id, base_price, start_date, end_date, is_active, priority
`

type UpdatePPriceParams struct {
	BasePrice sql.NullFloat64 `json:"base_price"`
	EndDate   sql.NullTime    `json:"end_date"`
	IsActive  bool            `json:"is_active"`
	Priority  sql.NullInt32   `json:"priority"`
	ID        int64           `json:"id"`
}

func (q *Queries) UpdatePPrice(ctx context.Context, arg UpdatePPriceParams) (ProductPricing, error) {
	row := q.db.QueryRowContext(ctx, updatePPrice,
		arg.BasePrice,
		arg.EndDate,
		arg.IsActive,
		arg.Priority,
		arg.ID,
	)
	var i ProductPricing
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.BasePrice,
		&i.StartDate,
		&i.EndDate,
		&i.IsActive,
		&i.Priority,
	)
	return i, err
}
