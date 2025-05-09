package main

import (
	"net/http"

	env "github.com/songphuc19102004/social/internal"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	responder := env.GetString("SERVICE_NAME", "invalid")
	env := env.GetString("ENV", "invalid")
	data := map[string]string{
		"status":    "ok",
		"responder": responder,
		"env":       env,
		"version":   version,
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
	}
}
