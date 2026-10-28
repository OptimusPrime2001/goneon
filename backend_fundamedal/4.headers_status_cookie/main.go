package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type errorResponse struct {
	Error string `json:"error"`
}

func findUser(id int) (user any, found bool) {
	return nil, true
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{"id must be a number"})
		return
	}

	user, found := findUser(id) // giả sử hàm này tìm trong DB
	if !found {
		writeJSON(w, http.StatusNotFound, errorResponse{"user not found"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func main() {
	cookieLogin()
}
