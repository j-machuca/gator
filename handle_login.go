package main

import (
	"context"
	"fmt"
)

func handleLogin(s *state, cmd command) error {
	ctx := context.Background()
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]
	err := s.cfg.SetUser(name)
	if err != nil {
		return err
	}

	_, err = s.db.GetUserByName(ctx, name)

	if err != nil {
		return fmt.Errorf("user does not exist %w", err)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Println("User switched succesfully!")
	return nil
}
