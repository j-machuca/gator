package main

import (
	"context"
	"fmt"

	"github.com/j-machuca/gator/internal/rss"
)

func handleAgg(s *state, cmd command) error {
	ctx := context.Background()
	rssFeed, err := rss.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")

	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}
	fmt.Printf("%+v\n", rssFeed)
	return nil
}
