package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/songphuc19102004/social/internal/store"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("userId")
	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := app.store.Users.GetById(ctx, userIdInt)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFound(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	err = app.jsonResponse(w, http.StatusOK, user)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
