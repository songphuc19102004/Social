package main

import (
	"log"
	"net/http"
	"time"

	"github.com/songphuc19102004/social/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() *http.ServeMux {
	mux := http.NewServeMux()

	// healthcheck
	mux.HandleFunc("GET /v1/health", app.healthCheckHandler)

	// posts
	mux.HandleFunc("POST /v1/posts", app.createPostHandler)
	mux.HandleFunc("GET /v1/posts/{postId}", app.getPostHandler)

	// users

	return mux
}

func (app *application) run(mux *http.ServeMux) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      Logger(mux),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("server is running on port %v", app.config.addr)
	return server.ListenAndServe()
}
