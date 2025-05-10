package main

import (
	"net/http"
	"strconv"

	"github.com/songphuc19102004/social/internal"
	"github.com/songphuc19102004/social/internal/store"
)

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,max=100"`
	UserId  int64  `json:"user_id"`
	PostId  int64  `json:"post_id"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var createRequest CreateCommentRequest
	err := readJSON(w, r, &createRequest)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := internal.Validate(createRequest); err != nil {
		app.badRequest(w, r, err)
		return
	}

	comment := &store.Comment{
		Content: createRequest.Content,
		PostID:  createRequest.PostId,
		// TODO: change this after authen
		UserID: 1,
	}

	ctx := r.Context()

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getCommentHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.PathValue("commentId")
	postIdInt, err := strconv.ParseInt(postId, 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()
	comment, err := app.store.Comments.GetById(ctx, postIdInt)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = writeJSON(w, http.StatusOK, comment)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
