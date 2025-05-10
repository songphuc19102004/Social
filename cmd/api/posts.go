package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/songphuc19102004/social/internal"
	"github.com/songphuc19102004/social/internal/store"
)

type key string

const (
	postKey key = "post"
)

type createPostPayload struct {
	Content string   `json:"content" validate:"min=2,max=100"`
	Title   string   `json:"title" validate:"min=2,max=100"`
	Tags    []string `json:"tags"`
}

type updatePostPayload struct {
	Content *string   `json:"content" validate:"min=2,max=100"`
	Title   *string   `json:"title" validate:"min=2,max=100"`
	Tags    *[]string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var createRequest createPostPayload

	if err := readJSON(w, r, &createRequest); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := internal.Validate(createRequest); err != nil {
		app.badRequest(w, r, err)
		return
	}

	post := &store.Post{
		Content: createRequest.Content,
		Title:   createRequest.Title,
		Tags:    createRequest.Tags,
		// change after auth
		UserID: 1,
	}
	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	ctx := r.Context()
	comments, err := app.store.Comments.GetByPostId(ctx, post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = *comments

	if err = writeJSON(w, http.StatusOK, post); err != nil {
		app.badRequest(w, r, err)
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	var updatePostPayload updatePostPayload
	if err := readJSON(w, r, &updatePostPayload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := internal.Validate(updatePostPayload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	post := getPostFromCtx(r)

	if updatePostPayload.Content != nil {
		post.Content = *updatePostPayload.Content
	}

	if updatePostPayload.Title != nil {
		post.Title = *updatePostPayload.Title
	}

	if updatePostPayload.Tags != nil {
		post.Tags = *updatePostPayload.Tags
	}

	ctx := r.Context()

	if err := app.store.Posts.Update(ctx, post); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, *post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	ctx := r.Context()
	err := app.store.Posts.Delete(ctx, post.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFound(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postId := r.PathValue("postId")
		postIdInt, err := strconv.ParseInt(postId, 10, 64)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetById(ctx, postIdInt)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFound(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, postKey, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	ctx := r.Context()
	post, _ := ctx.Value(postKey).(*store.Post)
	return post
}
