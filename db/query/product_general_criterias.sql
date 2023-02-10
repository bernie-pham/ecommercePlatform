-- name: CreatePCriteria :one
INSERT INTO 
    product_general_criteria (criteria)
VALUES ($1) RETURNING *;

-- name: DeletePCriteria :exec
DELETE FROM product_general_criteria
WHERE id = $1;

-- name: ListPCriterias :many
SELECT *
FROM product_general_criteria;