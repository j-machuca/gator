package main

import (
	"context"
	"fmt"
	"strconv"
)

func handleBrowse(s *state, cmd command) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("Usage: %s <limit>", cmd.Name)
	}
	var limit int32
	if len(cmd.Args) == 1 {
		v, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("failed to convert %s to integer", cmd.Args[0])
		}
		limit = int32(v)
	} else {
		limit = 2
	}
	posts, err := s.db.GetPosts(context.Background(), limit)
	if err != nil {
		return fmt.Errorf("Failed to retrieve posts Error:\n %w", err)
	}
	for _, p := range posts {
		fmt.Printf("Title: %v\n", p.Title)
		fmt.Printf("Description: %v\n", p.Description)
		fmt.Printf("Published Date: %v\n", p.PublishedAt)
		fmt.Printf("URL: %v\n", p.Url)
	}
	return nil
}
