package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(ctx context.Context, postId int64) (*Post, error)
		Delete(ctx context.Context, postId int64) error
		Update(ctx context.Context, post *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostId(ctx context.Context, postId int64) (*[]Comment, error)
		GetById(ctx context.Context, commentId int64) (*Comment, error)
	}
}

// We can create internal/postgres, internal/mongodb, internal/mysql
// if we have multiple db
func NewPostgresStorage(db *sql.DB) *Storage {
	return &Storage{
		Posts:    &PostStore{db},
		Users:    &UserStore{db},
		Comments: &CommentStore{db},
	}
}
