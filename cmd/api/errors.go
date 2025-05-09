package main

import (
	"log"
	"net/http"
	"time"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%v | internal server error path: %v, method: %v, error: %v\n", time.Now(), r.URL.Path, r.Method, err.Error())
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%v | bad request error path: %v, method: %v, error: %v\n", time.Now(), r.URL.Path, r.Method, err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%v | not found error path: %v, method: %v, error: %v\n", time.Now(), r.URL.Path, r.Method, err.Error())
	writeJSONError(w, http.StatusNotFound, "not found")
}
