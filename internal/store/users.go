package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID             int64
	Username       string
	FirstName      string
	LastName       string
	Age            int
	HashedPassword string
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
		&user.Username,
	)
	if err != nil {
		return err
	}

	return nil
}
