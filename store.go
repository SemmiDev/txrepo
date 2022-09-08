package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Store interface {
	Repository
	UpdateUserTx(ctx context.Context, userID uuid.UUID, user *User) error
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func (store *SQLStore) UpdateUserTx(ctx context.Context, userID uuid.UUID, user *User) error {
	err := store.execTx(ctx, func(q *Queries) error {
		exists, err := q.ExistsUser(ctx, userID)
		if err != nil {
			return err
		}
		if !exists {
			return ErrUserNotFound
		}
		err = q.UpdateUser(ctx, userID, user)
		return err
	})
	return err
}
