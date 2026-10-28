package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func Request() {
	client := &http.Client{Timeout: 10 * time.Second}

	body := strings.NewReader(`{"name":"An","email":"an@example.com"}`)
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost, "http://localhost:8080/v1/users", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer abc123")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode, resp.Header.Get("Location"))
	fmt.Println(string(data))
}
