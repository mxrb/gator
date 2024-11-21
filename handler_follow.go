package main

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/mxrb/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <url>", cmd.Name)
	}
	feedUrl := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("couldn't find feed: %w", err)
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}
	printFollowResult(feedFollow)
	return nil
}

func printFollowResult(feedFollow database.CreateFeedFollowRow) error {
	const templText = `Successfully followed feed:
Feed: {{.FeedName}}
User: {{.UserName}}
`
	t := template.Must(template.New("follow").Parse(templText))
	return t.Execute(os.Stdout, feedFollow)
}
