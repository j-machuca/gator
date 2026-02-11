package main

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"time"

	"github.com/j-machuca/gator/internal/database"
	"github.com/j-machuca/gator/internal/rss"
)

// need to refactor to call the scrapefeeds function
func handleAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <time_between_requests>", cmd.Name)
	}
	interval, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Failed to parse time between requests\n Error: %w", err)
	}
	fmt.Printf("Collecting feeds every %v\n", interval)
	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		fmt.Println("Starting the scrape loop")
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	fmt.Println("Fetching next feed to fetch")
	ctx := context.Background()
	feedToFetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("Failed to retrieve feed to fetch %w", err)
	}
	fmt.Printf("Marking %s as fetched\n ", feedToFetch.Url)
	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID: feedToFetch.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	})
	if err != nil {
		return fmt.Errorf("Failed to mark feed %v as fetched \n Error:%w", feedToFetch.ID, err)
	}
	fmt.Printf("Fetching data from %s\n", feedToFetch.Url)
	feed, err := rss.FetchFeed(ctx, feedToFetch.Url)
	if err != nil {
		return fmt.Errorf("Failed to scrape feed: %s error: %w", feedToFetch.Url, err)
	}

	fmt.Print("iterating over items in feed\n")
	for _, item := range feed.Channel.Item {
		fmt.Println(html.UnescapeString(item.Title))
	}

	return nil
}
