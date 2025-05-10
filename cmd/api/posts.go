package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/songphuc19102004/social/internal"
	"github.com/songphuc19102004/social/internal/store"
)

type createPostRequest struct {
	Content string   `json:"content" validate:"min=2,max=100"`
	Title   string   `json:"title" validate:"min=2,max=100"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var createRequest createPostRequest

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
	postId := r.PathValue("postId")

	// postIdInt, err := strconv.Atoi(postId)
	// if err != nil {
	// 	writeJSONError(w, http.StatusBadRequest, err)
	// 	return
	// }

	postIdInt, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	post, err := app.store.Posts.GetById(ctx, int64(postIdInt))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFound(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err = writeJSON(w, http.StatusOK, post); err != nil {
		app.badRequest(w, r, err)
	}
}

func updatePostHandler(w http.ResponseWriter, r *http.Request) {
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
}
