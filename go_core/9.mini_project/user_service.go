package main

import (
	"bufio"
	"errors"
	"fmt"
)

func createUser(scanner *bufio.Scanner, store *UserStore) {
	name, ok := readInput(scanner, "Name: ")
	if !ok {
		return
	}

	email, ok := readInput(scanner, "Email: ")
	if !ok {
		return
	}

	user, err := store.CreateUser(name, email)
	if err != nil {
		fmt.Println("Create user failed:", err)
		return
	}

	if err := store.Save(dataFile); err != nil {
		rollbackErr := store.DeleteUser(user.ID)
		fmt.Println("Save failed:", err)

		if rollbackErr != nil {
			fmt.Println("Rollback failed:", rollbackErr)
		}
		return
	}

	fmt.Printf("Created user: %+v\n", user)
}

func listUsers(store *UserStore) {
	users := store.GetAllUsers()

	if len(users) == 0 {
		fmt.Println("No users found.")
		return
	}

	for _, user := range users {
		fmt.Printf("ID: %d | Name: %s | Email: %s\n",
			user.ID, user.Name, user.Email)
	}
}

func getUserByID(scanner *bufio.Scanner, store *UserStore) {
	id, ok := readID(scanner)
	if !ok {
		return
	}

	user, err := store.GetUserByID(id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			fmt.Println("User not found.")
			return
		}

		fmt.Println("Get user failed:", err)
		return
	}

	fmt.Printf("ID: %d | Name: %s | Email: %s\n",
		user.ID, user.Name, user.Email)
}

func updateUser(scanner *bufio.Scanner, store *UserStore) {
	id, ok := readID(scanner)
	if !ok {
		return
	}

	if _, err := store.GetUserByID(id); err != nil {
		fmt.Println("Get user failed:", err)
		return
	}

	name, ok := readInput(scanner, "New name: ")
	if !ok {
		return
	}

	email, ok := readInput(scanner, "New email: ")
	if !ok {
		return
	}

	oldUser, _ := store.GetUserByID(id)

	if err := store.UpdateUser(id, name, email); err != nil {
		fmt.Println("Update user failed:", err)
		return
	}

	if err := store.Save(dataFile); err != nil {
		rollbackErr := store.UpdateUser(id, oldUser.Name, oldUser.Email)
		fmt.Println("Save failed:", err)

		if rollbackErr != nil {
			fmt.Println("Rollback failed:", rollbackErr)
		}
		return
	}

	fmt.Println("User updated successfully.")
}

func deleteUser(scanner *bufio.Scanner, store *UserStore) {
	id, ok := readID(scanner)
	if !ok {
		return
	}

	user, err := store.GetUserByID(id)
	if err != nil {
		fmt.Println("Get user failed:", err)
		return
	}

	if err := store.DeleteUser(id); err != nil {
		fmt.Println("Delete user failed:", err)
		return
	}

	if err := store.Save(dataFile); err != nil {
		_, rollbackErr := store.CreateUserWithID(user)
		fmt.Println("Save failed:", err)

		if rollbackErr != nil {
			fmt.Println("Rollback failed:", rollbackErr)
		}
		return
	}

	fmt.Println("User deleted successfully.")
}
