package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const dataFile = "users.json"

func main() {
	store, err := NewUserStore(dataFile)
	if err != nil {
		fmt.Println("Cannot initialize user store:", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n========== USER MANAGEMENT ==========")
		fmt.Println("1. Create user")
		fmt.Println("2. List users")
		fmt.Println("3. Get user by ID")
		fmt.Println("4. Update user")
		fmt.Println("5. Delete user")
		fmt.Println("0. Exit")
		fmt.Print("Choose an option: ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Println("Input error:", err)
			}
			return
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			createUser(scanner, store)

		case "2":
			listUsers(store)

		case "3":
			getUserByID(scanner, store)

		case "4":
			updateUser(scanner, store)

		case "5":
			deleteUser(scanner, store)

		case "0":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

func readInput(scanner *bufio.Scanner, prompt string) (string, bool) {
	fmt.Print(prompt)

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Println("Input error:", err)
		}
		return "", false
	}

	return strings.TrimSpace(scanner.Text()), true
}

func readID(scanner *bufio.Scanner) (int, bool) {
	input, ok := readInput(scanner, "Enter user ID: ")
	if !ok {
		return 0, false
	}

	id, err := strconv.Atoi(input)
	if err != nil || id < 1 {
		fmt.Println("ID must be a positive integer.")
		return 0, false
	}

	return id, true
}
