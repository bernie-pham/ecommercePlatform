// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: product_general_criterias.sql

package db

import (
	"context"
)

const createPCriteria = `-- name: CreatePCriteria :one
INSERT INTO 
    product_general_criteria (criteria)
VALUES ($1) RETURNING id, criteria
`

func (q *Queries) CreatePCriteria(ctx context.Context, criteria string) (ProductGeneralCriterium, error) {
	row := q.db.QueryRowContext(ctx, createPCriteria, criteria)
	var i ProductGeneralCriterium
	err := row.Scan(&i.ID, &i.Criteria)
	return i, err
}

const deletePCriteria = `-- name: DeletePCriteria :exec
DELETE FROM product_general_criteria
WHERE id = $1
`

func (q *Queries) DeletePCriteria(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePCriteria, id)
	return err
}

const listPCriterias = `-- name: ListPCriterias :many
SELECT id, criteria
FROM product_general_criteria
`

func (q *Queries) ListPCriterias(ctx context.Context) ([]ProductGeneralCriterium, error) {
	rows, err := q.db.QueryContext(ctx, listPCriterias)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ProductGeneralCriterium{}
	for rows.Next() {
		var i ProductGeneralCriterium
		if err := rows.Scan(&i.ID, &i.Criteria); err != nil {
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
