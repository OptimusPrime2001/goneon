package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Xin chào!")
}

func type_handle() {
	mux := http.NewServeMux()
	mux.Handle("GET /hello", http.HandlerFunc(hello)) // ép kiểu hàm thành Handler
	mux.HandleFunc("GET /hello", hello)
	http.ListenAndServe(":8080", mux)
}

func main() {
	structHandler()
}
