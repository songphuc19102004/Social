package store

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	ID             int64
	Username       string
	FirstName      string
	LastName       string
	Age            int
	HashedPassword string
	CreatedAt      string
}

type UserStore struct {
	db *sql.DB
}

func (store *UserStore) Create(ctx context.Context, user *User) error {
	query := `
    INSERT INTO
    users(username, first_name, last_name, age, hashed_password)
    VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := store.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Age,
		user.HashedPassword,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (store *UserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	var user User
	query := `
  SELECT u.id, u.username, u.first_name, u.last_name, u.age, u.created_at
  FROM users u
  WHERE u.id = $1
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := store.db.QueryRowContext(ctx, query, userId).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
