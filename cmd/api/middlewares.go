package main

import (
	"log"
	"net/http"
	"time"

	env "github.com/songphuc19102004/social/internal"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Println(time.Since(start), r.Method, r.URL.Path, env.GetString("SERVICE_NAME", "env not found"))
	})
}
