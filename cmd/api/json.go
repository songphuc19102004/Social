package main

import (
	"encoding/json"
	"net/http"
)

const maxRequestBytes int64 = 1_048_576 // 1MB

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data *any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, err error) error {
	type errorMessage struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, errorMessage{Error: err.Error()})
}
