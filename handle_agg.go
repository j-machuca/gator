package main

import (
	"context"
	"database/sql"
	"fmt"
	"html"
	"time"

	"github.com/google/uuid"
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
		insertPost(s, item, feedToFetch.ID)
	}

	return nil
}

func insertPost(s *state, item rss.RSSItem, feed_id uuid.UUID) (database.Post, error) {
	p, err := s.db.InsertPost(context.Background(), database.InsertPostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		Title:       html.UnescapeString(item.Title),
		Description: html.UnescapeString(item.Description),
		Url:         item.Link,
		PublishedAt: item.PubDate,
		FeedID:      feed_id,
	})
	if err != nil {
		return database.Post{}, fmt.Errorf("Failed to save post to Database Error: \n %w", err)
	}
	fmt.Println("Inserted Post:")
	fmt.Printf("Title: %s\n", p.Title)
	fmt.Printf("Description: %s\n", p.Description)
	fmt.Printf("URL: %s\n", p.Url)
	return p, nil
}
