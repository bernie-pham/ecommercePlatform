-- name: CreatePColour :one
INSERT INTO 
    product_colour (colour_name)
VALUES ($1) RETURNING *;

-- name: DeletePColour :exec
DELETE FROM product_colour
WHERE id = $1;

-- name: ListPColours :many
SELECT *
FROM product_colour;