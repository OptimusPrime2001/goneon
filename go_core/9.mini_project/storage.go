package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func NewUserStore(path string) (*UserStore, error) {
	store := &UserStore{
		users:  []User{},
		nextID: 1,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return store, nil
		}

		return nil, fmt.Errorf("read users file %q: %w", path, err)
	}

	if err := json.Unmarshal(data, &store.users); err != nil {
		return nil, fmt.Errorf("decode users file %q: %w", path, err)
	}

	if store.users == nil {
		store.users = []User{}
	}

	for _, user := range store.users {
		if user.ID < 1 {
			return nil, fmt.Errorf("invalid user ID: %d", user.ID)
		}

		if user.ID >= store.nextID {
			store.nextID = user.ID + 1
		}
	}

	for i := range store.users {
		for j := i + 1; j < len(store.users); j++ {
			if store.users[i].ID == store.users[j].ID {
				return nil, fmt.Errorf("%w: ID=%d", ErrDuplicateID, store.users[i].ID)
			}
		}
	}

	return store, nil
}

func (s *UserStore) Save(path string) error {
	data, err := json.MarshalIndent(s.users, "", "    ")
	if err != nil {
		return fmt.Errorf("encode users: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write users file %q: %w", path, err)
	}

	return nil
}
