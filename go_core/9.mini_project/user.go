package main

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrDuplicateID  = errors.New("user ID already exists")
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserStore struct {
	users  []User
	nextID int
}

func (s *UserStore) CreateUser(name, email string) (User, error) {
	if name == "" || email == "" {
		return User{}, fmt.Errorf("%w: name and email are required", ErrInvalidInput)
	}

	user := User{
		ID:    s.nextID,
		Name:  name,
		Email: email,
	}

	s.users = append(s.users, user)
	s.nextID++

	return user, nil
}

func (s *UserStore) GetAllUsers() []User {
	users := make([]User, len(s.users))
	copy(users, s.users)

	return users
}

func (s *UserStore) GetUserByID(id int) (User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("%w: ID=%d", ErrUserNotFound, id)
}

func (s *UserStore) UpdateUser(id int, name, email string) error {
	if name == "" || email == "" {
		return fmt.Errorf("%w: name and email are required", ErrInvalidInput)
	}

	for i := range s.users {
		if s.users[i].ID == id {
			s.users[i].Name = name
			s.users[i].Email = email

			return nil
		}
	}

	return fmt.Errorf("%w: ID=%d", ErrUserNotFound, id)
}

func (s *UserStore) DeleteUser(id int) error {
	for i := range s.users {
		if s.users[i].ID == id {
			s.users = append(s.users[:i], s.users[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("%w: ID=%d", ErrUserNotFound, id)
}
func (s *UserStore) CreateUserWithID(user User) (User, error) {
	if user.ID < 1 || user.Name == "" || user.Email == "" {
		return User{}, fmt.Errorf("%w: invalid user", ErrInvalidInput)
	}

	for _, existing := range s.users {
		if existing.ID == user.ID {
			return User{}, fmt.Errorf("%w: ID=%d", ErrDuplicateID, user.ID)
		}
	}

	s.users = append(s.users, user)

	if user.ID >= s.nextID {
		s.nextID = user.ID + 1
	}

	return user, nil
}
