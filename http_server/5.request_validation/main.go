package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorBody struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

func writeError(w http.ResponseWriter, status int, code, msg string, details ...FieldError) {
	writeJSON(w, status, map[string]ErrorBody{
		"error": {Code: code, Message: msg, Details: details},
	})
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserStore struct {
	mu     sync.Mutex
	nextID int
	users  []User
}

type UserHandler struct{ store *UserStore }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// ... FieldError, ErrorBody, writeError, decodeJSON, CreateUserRequest + Validate,
//
//	isEmail, ErrEmailTaken, UserStore.Create, UserHandler.Create ở các mục trên ...
var ErrEmailTaken = errors.New("email already taken")

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if !decodeJSON(w, r, &req) { // bước 1: decode
		return
	}
	if errs := req.Validate(); len(errs) > 0 { // bước 2: validate field
		writeError(w, http.StatusUnprocessableEntity, "validation_failed", "Dữ liệu không hợp lệ", errs...)
		return
	}

	user, err := h.store.Create(req.Name, req.Email) // bước 3: nghiệp vụ
	switch {
	case errors.Is(err, ErrEmailTaken):
		writeError(w, http.StatusConflict, "email_taken", "email is already in use")
		return
	case err != nil:
		slog.Error("create user failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
		return
	}

	w.Header().Set("Location", "/users/"+strconv.Itoa(user.ID))
	writeJSON(w, http.StatusCreated, user)
}
func (s *UserStore) Create(name, email string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.users {
		if u.Email == email {
			return User{}, ErrEmailTaken
		}
	}
	s.nextID++
	u := User{ID: s.nextID, Name: name, Email: email}
	s.users = append(s.users, u)
	return u, nil
}
func main() {
	h := &UserHandler{store: &UserStore{}}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", h.Create)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	slog.Info("listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
