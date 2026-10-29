package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// Handler dùng chung: in ra thông tin mux đã khớp được
func show(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"method":  r.Method,
		"path":    r.URL.Path,
		"pattern": r.Pattern,
		"id":      r.PathValue("id"),
		"rest":    r.PathValue("path"),
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", show)
	mux.HandleFunc("GET /users", show)
	mux.HandleFunc("POST /users", show)
	mux.HandleFunc("GET /users/{id}", show)
	mux.HandleFunc("GET /users/me", show)
	mux.HandleFunc("DELETE /users/{id}", show)
	mux.HandleFunc("GET /files/{path...}", show)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
