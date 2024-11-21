package main

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/mxrb/gator/internal/database"
)

func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %v <name> <url>", cmd.Name)
	}
	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{Name: feedName, Url: feedUrl, UserID: user.ID})
	if err != nil {
		return fmt.Errorf("couldn't save feed: %w", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}
	fmt.Println("Feed created successfully:")
	if err := printFeed(feed); err != nil {
		return fmt.Errorf("couldn't print out feed: %w", err)
	}
	fmt.Println("===============================================")
	return nil
}

func printFeed(feed database.Feed) error {
	const templText = `ID:      {{.ID}}
Name:    {{.Name}}
URL:     {{.Url}}
UserID:  {{.UserID}}
Created: {{.CreatedAt}}
Updated: {{.UpdatedAt}}
`
	t := template.Must(template.New("feed").Parse(templText))
	return t.Execute(os.Stdout, feed)
}
