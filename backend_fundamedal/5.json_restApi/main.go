package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Password    string    `string:"password"`
	Bio         string    `json:"bio,omitempty"`
	Age         *int      `json:"age"`
	CreatedAt   time.Time `json:"created_at"`
	Price_cents int64     `json:"price_cents"`
}

func main() {
	user := User{
		ID:        2001,
		Name:      "TrungLe",
		Password:  "Trung2001@",
		CreatedAt: time.Now(),
	}
	data, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))

	var user2 User
	input := `{"id": 2, "name": "Bình", "age": 25}`
	if err := json.Unmarshal([]byte(input), &user2); err != nil {
		panic(err)
	}
	fmt.Println(user2.Name, *user2.Age)
}
