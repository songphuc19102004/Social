package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var ErrNotFound = errors.New("record not found")

type Post struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title" validate:"required,max=100"`
	Tags      []string `json:"tags"`
	UserID    int      `json:"user_id"`
	CreatedAt string   `json:"created_at"`
	Updatedt  string   `json:"updated_at"`
}

type PostStore struct {
	db *sql.DB
}

func (store *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
    INSERT INTO
    posts(content, title, tags, user_id)
    VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at
  `

	err := store.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		pq.Array(post.Tags),
		post.UserID,
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.Updatedt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (store *PostStore) GetById(ctx context.Context, postId int64) (*Post, error) {
	query := `
    SELECT id, content, title, tags, user_id, created_at, updated_at
    FROM posts p
    WHERE p.id = $1
  `

	var post Post

	err := store.db.QueryRowContext(ctx, query, postId).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		pq.Array(&post.Tags),
		&post.UserID,
		&post.CreatedAt,
		&post.Updatedt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}
