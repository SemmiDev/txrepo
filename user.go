package main

import (
	"context"

	"github.com/google/uuid"
)

type Creator interface {
	CreateUser(ctx context.Context, user *User) error
}

type Updater interface {
	UpdateUser(ctx context.Context, userID uuid.UUID, user *User) error
}

type Finder interface {
	FindUser(ctx context.Context, userID uuid.UUID) (*User, error)
	FindUsers(ctx context.Context) ([]*User, error)
}

type User struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func NewUser(name, email string) *User {
	return &User{
		ID:    uuid.New(),
		Name:  name,
		Email: email,
	}
}

func NewUserWithID(id uuid.UUID, name, email string) *User {
	return &User{
		ID:    id,
		Name:  name,
		Email: email,
	}
}

func NewUserWithoutID(name, email string) *User {
	return &User{
		Name:  name,
		Email: email,
	}
}

const createUserQuery = `INSERT INTO users (id,name,email) VALUES ($1,$2,$3)`

func (q *Queries) CreateUser(ctx context.Context, user *User) error {
	_, err := q.db.ExecContext(ctx, createUserQuery,
		user.ID,
		user.Name,
		user.Email,
	)
	return err
}

const findUserQuery = `SELECT id,name,email FROM users WHERE id=$1`

func (q *Queries) FindUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	row := q.db.QueryRowContext(ctx, findUserQuery, userID.String())
	var u User
	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

const findUsersQuery = `SELECT id,name,email FROM users`

func (q *Queries) FindUsers(ctx context.Context) ([]*User, error) {
	rows, err := q.db.QueryContext(ctx, findUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User
	for rows.Next() {
		var u User
		err = rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

const existsUserQuery = `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`

func (q *Queries) ExistsUser(ctx context.Context, userID uuid.UUID) (bool, error) {
	row := q.db.QueryRowContext(ctx, existsUserQuery, userID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const updateUserQuery = `UPDATE users SET name=$1,email=$2 WHERE id=$3`

func (q *Queries) UpdateUser(ctx context.Context, userID uuid.UUID, user *User) error {
	_, err := q.db.ExecContext(ctx, updateUserQuery,
		user.Name,
		user.Email,
		userID,
	)
	return err
}
