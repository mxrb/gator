package main

import (
	"context"
	"fmt"

	"github.com/mxrb/gator/internal/database"
)

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <url>", cmd.Name)
	}
	feedUrl := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("couldn't find feed: %w", err)
	}
	changed, err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return fmt.Errorf("couldn't delete feed follow: %w", err)
	}
	if changed == 0 {
		fmt.Println("User doesn't follow this feed.")
	}
	fmt.Println("Successfully unfollowed feed.")
	return nil
}
