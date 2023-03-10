// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: notification.sql

package db

import (
	"context"
)

const createNotification = `-- name: CreateNotification :one
INSERT INTO notifications (
    message,
    recipient_id,
    title
)
VALUES ($1, $2, $3)
RETURNING id, message, recipient_id, created_at, title
`

type CreateNotificationParams struct {
	Message     string `json:"message"`
	RecipientID int64  `json:"recipient_id"`
	Title       string `json:"title"`
}

func (q *Queries) CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error) {
	row := q.db.QueryRowContext(ctx, createNotification, arg.Message, arg.RecipientID, arg.Title)
	var i Notification
	err := row.Scan(
		&i.ID,
		&i.Message,
		&i.RecipientID,
		&i.CreatedAt,
		&i.Title,
	)
	return i, err
}

const listNotifications = `-- name: ListNotifications :many
SELECT id, message, recipient_id, created_at, title
FROM notifications
WHERE recipient_id = $1
`

func (q *Queries) ListNotifications(ctx context.Context, recipientID int64) ([]Notification, error) {
	rows, err := q.db.QueryContext(ctx, listNotifications, recipientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Notification{}
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.Message,
			&i.RecipientID,
			&i.CreatedAt,
			&i.Title,
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
