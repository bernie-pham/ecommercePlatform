-- name: CreateNotification :one
INSERT INTO notifications (
    message,
    recipient_id,
    title
)
VALUES ($1, $2, $3)
RETURNING *;


-- name: ListNotifications :many
SELECT *
FROM notifications
WHERE recipient_id = $1;
