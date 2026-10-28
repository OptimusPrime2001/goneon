package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func BAI7() {
	file, error := os.Open("users.json")
	if error != nil {
		fmt.Println("Open file error:", error)
		return
	}
	defer file.Close()

	var user []User
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&user); err != nil {
		fmt.Println("Decode error:", err)
		return
	}
	fmt.Printf("%+v\n", user)
}
func BAI7_Write() {
	file, err := os.Create("output.json")
	if err != nil {
		fmt.Println("Create file error:", err)
		return
	}
	defer file.Close()

	user := User{
		ID:    2,
		Name:  "Bob",
		Email: "bob@example.com",
	}

	encoder := json.NewEncoder(file) //tạo encoder ghi JSON trực tiếp vào file.
	encoder.SetIndent("", "    ")    //giúp JSON dễ đọc.

	if err := encoder.Encode(user); err != nil { //ghi JSON và tự thêm ký tự xuống dòng ở cuối.
		fmt.Println("Encode error:", err)
		return
	}

	fmt.Println("User saved successfully")
}
func BAI7Stream() {
	file, err := os.Open("users.ndjson")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("File không tồn tại")
			return
		}

		fmt.Println("Open file error:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		var user User

		err := decoder.Decode(&user)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			fmt.Println("Decode error:", err)
			return
		}

		fmt.Printf("ID: %d, Name: %s\n", user.ID, user.Name)
	}
}
func main() {
	BAI7Stream()
}
