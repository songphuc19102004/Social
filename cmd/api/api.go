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
	mainMux := http.NewServeMux()
	v1Router := http.NewServeMux()
	mainMux.Handle("/v1/", http.StripPrefix("/v1", v1Router))

	// healthcheck
	v1Router.HandleFunc("GET /health", app.healthCheckHandler)

	// posts
	v1Router.HandleFunc("POST /posts", app.createPostHandler)

	// handling endpoints that requires getting a post from id
	postDetailRouter := http.NewServeMux()
	v1Router.Handle("/posts/{postId}", app.postContextMiddleware(postDetailRouter))
	v1Router.Handle("/posts/{postId}/", app.postContextMiddleware(postDetailRouter))

	postDetailRouter.HandleFunc("GET /", app.getPostHandler)
	postDetailRouter.HandleFunc("DELETE /", app.deletePostHandler)
	postDetailRouter.HandleFunc("PUT /", app.updatePostHandler)

	// users

	// comments
	v1Router.HandleFunc("POST /comments", app.createCommentHandler)
	v1Router.HandleFunc("GET /comments/{commentId}", app.getCommentHandler)

	return mainMux
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
