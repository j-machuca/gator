package main

import (
	"context"
	"fmt"
)

func handleReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("reset takes no arguments")
	}

	err := s.cfg.SetUser("")
	if err != nil {
		return err
	}
	err = s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}
	err = s.db.ResetFeed(context.Background())
	if err != nil {
		return err
	}
	err = s.db.ResetFeedFollowings(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("reset was successful!")
	return nil
}
