package main

import (
	"fmt"
	"net/http"

	env "github.com/songphuc19102004/social/internal"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Health ok from cluster %v", env.GetString("SERVICE_NAME", "env not found"))
	w.Write([]byte(msg))
}
