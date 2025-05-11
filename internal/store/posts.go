package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title" validate:"required,max=100"`
	Tags      []string  `json:"tags"`
	UserID    int       `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	Updatedt  string    `json:"updated_at"`
	User      User      `json:"user"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
    SELECT id, content, title, tags, user_id, created_at, updated_at, version
    FROM posts p
    WHERE p.id = $1
  `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post

	err := store.db.QueryRowContext(ctx, query, postId).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		pq.Array(&post.Tags),
		&post.UserID,
		&post.CreatedAt,
		&post.Updatedt,
		&post.Version,
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

func (store *PostStore) Delete(ctx context.Context, postId int64) error {
	query := `
  DELETE FROM posts p
  WHERE p.id = $1;
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := store.db.Exec(query, postId)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (store *PostStore) Update(ctx context.Context, post *Post) error {
	query := `
  UPDATE posts
  SET title = $1, content = $2, tags = $3, version = version + 1
  WHERE id = $4 AND version = $5
  RETURNING version
  `

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := store.db.QueryRowContext(ctx, query, post.Title, post.Content, pq.Array(post.Tags), post.ID, post.Version).Scan(
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}
