package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/j-machuca/gator/internal/database"
)

func handleAddFeed(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("Failed to follow feed after creating %w", err)
	}
	_, err = s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("Failed to add a feed %w", err)
	}
	fmt.Println("Feed:")
	fmt.Printf(" * ID:%s\n", feed.ID)
	fmt.Printf(" * Name:%s\n", feed.Name)
	fmt.Printf(" * URL:%s\n", feed.Url)
	fmt.Printf(" * UserId:%s\n", feed.UserID)
	return nil
}

func handleGetFeeds(s *state, cmd command) error {
	ctx := context.Background()
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v", cmd.Name)
	}
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("Couldn't get feeds %w", err)
	}
	for _, feed := range feeds {
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(feed.Username.String)
	}
	return nil

}

func handleFeedFollow(s *state, cmd command, u database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]
	f, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error fetching feed by url: %w", err)
	}
	ff, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    u.ID,
		FeedID:    f.ID,
	})
	if err != nil {
		return fmt.Errorf("Error creating a feed follow %w", err)
	}

	fmt.Printf("Feed Name %s \n", ff.Feed)
	fmt.Printf("Current Username %s \n", u.Name)
	return nil

}

func handleGetFeedsForUser(s *state, cmd command, u database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("Usage: %s", cmd.Name)
	}
	ff, err := s.db.GetFeedFollowsForUser(context.Background(), u.ID)
	if err != nil {
		return fmt.Errorf("Failed to retrieve feed following for user %w", err)
	}
	for _, f := range ff {
		fmt.Println(f.Feed)
	}
	return nil
}

func handleUnfollow(s *state, cmd command, u database.User) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: unfollow <url>")
	}

	err := s.db.Unfollow(context.Background(), database.UnfollowParams{
		Url:  cmd.Args[0],
		Name: u.Name,
	})
	if err != nil {
		return fmt.Errorf("Failed unfollowing feed %w", err)
	}
	return nil
}
