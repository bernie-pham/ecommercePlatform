package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	ResetPasswdTX(ctx context.Context, arg ResetPasswdTXParams) error
	// CreateOrderTX(ctx context.Context, arg CreateOrderTXParams) error
	CreateOrderTX_v2(ctx context.Context, arg CreateOrderTXParams) error
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SQLStore) execTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
