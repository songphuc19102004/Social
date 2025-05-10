package store

import (
	"context"
	"database/sql"
	"errors"
)

type Comment struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	User      User   `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func (store *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
  INSERT INTO comments(content, post_id, user_id)
  VALUES ($1, $2, $3) RETURNING id, created_at, updated_at;
  `
	err := store.db.QueryRowContext(ctx, query, comment.Content, comment.PostID, comment.UserID).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (store *CommentStore) GetById(ctx context.Context, commentId int64) (*Comment, error) {
	query := `
  SELECT c.id, c.content, c.post_id, c.user_id, c.created_at, c.updated_at
  FROM comments c
  WHERE c.id = $1
  `
	var comment Comment
	err := store.db.QueryRowContext(ctx, query, commentId).Scan(
		&comment.ID,
		&comment.Content,
		&comment.PostID,
		&comment.UserID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &comment, err
}

func (store *CommentStore) GetByPostId(ctx context.Context, postId int64) (*[]Comment, error) {
	query := `
  SELECT c.id, c.content, c.post_id, c.user_id, u.username, c.created_at, c.updated_at
  FROM comments c
  JOIN users u ON c.user_id = u.id
  WHERE c.post_id = $1
  ORDER BY c.created_at DESC;
  `

	rows, err := store.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	comments := []Comment{}
	for rows.Next() {
		var c Comment
		err := rows.Scan(
			&c.ID,
			&c.Content,
			&c.PostID,
			&c.UserID,
			&c.User.Username,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return &comments, nil
}
