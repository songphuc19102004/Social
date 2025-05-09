package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(ctx context.Context, postId int64) (*Post, error)
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

// We can create internal/postgres, internal/mongodb, internal/mysql
// if we have multiple db
func NewPostgresStorage(db *sql.DB) *Storage {
	return &Storage{
		Posts: &PostStore{db},
		Users: &UserStore{db},
	}
}
