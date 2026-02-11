package main

import (
	"context"
	"fmt"

	"github.com/j-machuca/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return fmt.Errorf("Failed to retrieve user %w", err)
		}
		return handler(s, cmd, user)
	}
}
