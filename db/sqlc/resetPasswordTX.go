package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ResetPasswdTXParams struct {
	Email             string    `json:"email"`
	VerificationID    uuid.UUID `json:"verification_id"`
	NewHashedPassword string    `json:"hashed_password"`
}

func (store *SQLStore) ResetPasswdTX(ctx context.Context, arg ResetPasswdTXParams) error {
	err := store.execTX(ctx, func(q *Queries) error {
		err := q.OccupyVerification(ctx, arg.VerificationID)
		if err != nil {
			return err
		}

		_, err = q.UpdateUser(ctx, UpdateUserParams{
			Email: arg.Email,
			HashedPassword: sql.NullString{
				String: arg.NewHashedPassword,
				Valid:  true,
			},
			PasswordUpdatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
