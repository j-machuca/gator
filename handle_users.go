package main

import (
	"context"
	"fmt"
)

func handleUsers(s *state, cmd command) error {
	ctx := context.Background()
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: users")
	}

	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUsername {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}
