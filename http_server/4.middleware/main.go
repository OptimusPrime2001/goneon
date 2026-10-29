package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// ... các middleware và helper ở trên ...

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) // dev có thể dùng NewTextHandler

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"id": r.PathValue("id")})
	})
	mux.HandleFunc("GET /panic", func(w http.ResponseWriter, r *http.Request) {
		panic("boom") // để thử middleware recover
	})

	handler := Chain(mux,
		recoverer(logger),
		requestID,
		logging(logger),
		cors([]string{"http://localhost:3000"}),
	)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	logger.Info("listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
