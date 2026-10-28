package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
)

var (
	mu       sync.Mutex
	sessions = map[string]string{} // session_id -> username
)

func newSessionID() string {
	b := make([]byte, 16)
	rand.Read(b) // crypto/rand: ngẫu nhiên đủ an toàn, không dùng math/rand
	return hex.EncodeToString(b)
}

func login(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Username == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username is required"})
		return
	}

	id := newSessionID()
	mu.Lock()
	sessions[id] = in.Username
	mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    id,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false, // đặt true khi chạy HTTPS ở production
		SameSite: http.SameSiteLaxMode,
	})
	writeJSON(w, http.StatusOK, map[string]string{"message": "logged in"})
}

func me(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_id") // trả http.ErrNoCookie nếu không có
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "not logged in"})
		return
	}

	mu.Lock()
	user, ok := sessions[c.Value]
	mu.Unlock()
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid session"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"username": user})
}

func logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("session_id"); err == nil {
		mu.Lock()
		delete(sessions, c.Value) // hủy session ở phía server
		mu.Unlock()
	}
	// Xóa cookie ở trình duyệt: đặt MaxAge âm
	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
	w.WriteHeader(http.StatusNoContent)
}

func cookieLogin() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", login)
	mux.HandleFunc("GET /me", me)
	mux.HandleFunc("POST /logout", logout)
	http.ListenAndServe(":8080", mux)
}
