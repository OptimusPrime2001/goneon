package main

import (
	"encoding/json"
	"net/http"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// 1. Chặn body quá lớn (ở đây 1 MB), tránh bị làm quá tải bộ nhớ
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	// 2. Đọc và parse JSON
	var req CreateUserRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // từ chối field lạ để bắt lỗi gõ sai sớm
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return // đừng quên!
	}

	// 3. Xử lý nghiệp vụ (ở đây giả lập)
	user := UserResponse{ID: 42, Name: req.Name, Email: req.Email}

	// 4. Trả response: header -> status -> body
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/v1/users/42")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(user)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/users", createUser)
	http.ListenAndServe(":8080", mux)
}
