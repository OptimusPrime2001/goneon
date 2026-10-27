package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive bool   `json:"isActive"`
}

func BAI3() {
	user := User{
		ID:       1,
		Name:     "Trung Le",
		Email:    "Trungthekop123@gmail.com",
		IsActive: false,
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}
	fmt.Println(string(data))

	rawJSON := []byte(`{
        "id": 2,
        "name": "Alice",
        "email": "alice@example.com",
        
				"is_active": true
    }`)
	var user2 User

	err2 := json.Unmarshal(rawJSON, &user2)
	if err2 != nil {
		fmt.Println("Unmarshal error:", err)
		return
	}
	data, err3 := json.MarshalIndent(user2, "", "    ")
	if err3 != nil {
		fmt.Println("Marshal error:", err)
		return
	}
	fmt.Printf("\n%+v\n", user2)

	fmt.Println(string(data))
}
func BAI4() {
	// 1. Đọc file
	file, err := os.ReadFile("users.json")
	if err != nil {
		fmt.Println("Không tồn tài file users.json")
	}

	// 2. Chuyển JSON thành slice User
	var users []User
	if err := json.Unmarshal(file, &users); err != nil {
		fmt.Println("Parse json bị lỗi", err)
	}

	// 3. Cập nhật email của user có ID = 2
	for i, user := range users {
		if user.ID == 2 {
			users[i].Email = "trungthekop123@gmail.com"
		}
	}

	// 4. Chuyển slice User thành JSON
	updatedData, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}

	// 5. Ghi dữ liệu mới vào file
	if err := os.WriteFile("users.json", updatedData, 0644); err != nil {
		fmt.Println("Write file error:", err)
		return
	}
	fmt.Println("Users updated successfully")
}
func BAI5() {
	file, err := os.Open("large.txt")
	if err != nil {
		fmt.Println("Open file error:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Read file error:", err)
		return
	}
}
func BAI6() {
	file, err := os.Open("large.txt")
	if err != nil {
		fmt.Println("Có lỗi khi mở file: ", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totaLine, errorLines := 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), "error") {
			fmt.Println(line)
			errorLines += 1
		}
		totaLine += 1
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Đọc file gặp lỗi : ", err)
	}
	fmt.Println("\nTotal lines: ", totaLine)
	fmt.Println("Error lines: ", errorLines)
}
func main() {
	BAI6()
}
