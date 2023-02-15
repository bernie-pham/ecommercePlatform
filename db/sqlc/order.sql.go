// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: order.sql

package db

import (
	"context"
	"database/sql"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (user_id, base_price, discount_price)
VALUES ($1, $2, $3) RETURNING id
`

type CreateOrderParams struct {
	UserID        int64   `json:"user_id"`
	BasePrice     float32 `json:"base_price"`
	DiscountPrice float32 `json:"discount_price"`
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, createOrder, arg.UserID, arg.BasePrice, arg.DiscountPrice)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createOrderV2 = `-- name: CreateOrderV2 :one
INSERT INTO orders (
    user_id, 
    base_price, 
    discount_price,
    deal_id)
VALUES (
    $1, 
    $2, 
    $3, 
    $4
) RETURNING id
`

type CreateOrderV2Params struct {
	UserID        int64         `json:"user_id"`
	BasePrice     float32       `json:"base_price"`
	DiscountPrice float32       `json:"discount_price"`
	DealID        sql.NullInt64 `json:"deal_id"`
}

func (q *Queries) CreateOrderV2(ctx context.Context, arg CreateOrderV2Params) (int64, error) {
	row := q.db.QueryRowContext(ctx, createOrderV2,
		arg.UserID,
		arg.BasePrice,
		arg.DiscountPrice,
		arg.DealID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getCurrentOrder = `-- name: GetCurrentOrder :one
SELECT id, user_id, status, created_at, deal_id, base_price, discount_price
FROM orders
WHERE user_id = $1 AND status = 'open'
`

func (q *Queries) GetCurrentOrder(ctx context.Context, userID int64) (Order, error) {
	row := q.db.QueryRowContext(ctx, getCurrentOrder, userID)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Status,
		&i.CreatedAt,
		&i.DealID,
		&i.BasePrice,
		&i.DiscountPrice,
	)
	return i, err
}

const listOrder = `-- name: ListOrder :many
SELECT id, user_id, status, created_at, deal_id, base_price, discount_price 
FROM orders 
WHERE user_id = $1
`

func (q *Queries) ListOrder(ctx context.Context, userID int64) ([]Order, error) {
	rows, err := q.db.QueryContext(ctx, listOrder, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Order{}
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Status,
			&i.CreatedAt,
			&i.DealID,
			&i.BasePrice,
			&i.DiscountPrice,
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

const updateOrder = `-- name: UpdateOrder :one
UPDATE orders 
SET 
    status = COALESCE($1, status),
    deal_id = COALESCE($2, deal_id),
    base_price = COALESCE($3, base_price),
    discount_price = COALESCE($4, discount_price)
WHERE
    id = $5
RETURNING id, user_id, status, created_at, deal_id, base_price, discount_price
`

type UpdateOrderParams struct {
	Status        NullOrderStatus `json:"status"`
	DealID        sql.NullInt64   `json:"deal_id"`
	BasePrice     sql.NullFloat64 `json:"base_price"`
	DiscountPrice sql.NullFloat64 `json:"discount_price"`
	ID            int64           `json:"id"`
}

func (q *Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, updateOrder,
		arg.Status,
		arg.DealID,
		arg.BasePrice,
		arg.DiscountPrice,
		arg.ID,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Status,
		&i.CreatedAt,
		&i.DealID,
		&i.BasePrice,
		&i.DiscountPrice,
	)
	return i, err
}

const updateOrderStatus = `-- name: UpdateOrderStatus :exec
UPDATE orders 
SET 
    status = $1
WHERE id = $2
`

type UpdateOrderStatusParams struct {
	Status OrderStatus `json:"status"`
	ID     int64       `json:"id"`
}

func (q *Queries) UpdateOrderStatus(ctx context.Context, arg UpdateOrderStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateOrderStatus, arg.Status, arg.ID)
	return err
}
