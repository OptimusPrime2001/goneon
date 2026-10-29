package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("encode error:", err)
	}
}

// ---- Helpers ----

type errorBody struct {
	Error string `json:"error"`
}
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type UserStore struct {
	mu    sync.RWMutex
	users []User
}

func NewUserStore() *UserStore {
	return &UserStore{users: []User{
		{
			Id:   1,
			Name: "An",
		}, {
			Id:   2,
			Name: "Bình",
		},
	}}
}
func (store *UserStore) GetListUser() []User {
	store.mu.RLock()
	defer store.mu.RUnlock()
	result := make([]User, len(store.users))
	copy(result, store.users)
	return result
}

func (store *UserStore) GetUser(id int) (User, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	for _, user := range store.users {
		if user.Id == id {
			return user, true
		}
	}
	return User{}, false
}

type UserHandler struct {
	store *UserStore
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.GetListUser())
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue("id")
	id, err := strconv.Atoi(rawId)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "id must be a number"})
		return
	}
	user, ok := h.store.GetUser(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, errorBody{Error: "user not found"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func structHandler() {
	uh := &UserHandler{store: NewUserStore()}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", uh.List)
	mux.HandleFunc("GET /users/{id}", uh.Get)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Println("listening on", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

}
