package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/mxrb/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("argument is not a valid duration: %w", err)
	}
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(context.Background(), s)
	}
}

func scrapeFeeds(ctx context.Context, s *state) {
	nextDbFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Fatalf("couldn't get feed from database: %s", err)
	}
	scrapeFeed(ctx, s.db, nextDbFeed)
}

func scrapeFeed(ctx context.Context, db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		log.Printf("couldn't mark feed as fetched: %s", err)
		return
	}
	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		log.Printf("couldn't fetch feed: %s", err)
		return
	}
	for _, item := range rssFeed.Channel.Item {
		publishedAt, err := parsePubDate(item.PubDate)
		if err != nil {
			log.Printf("couldn't parse the publishing date: %s", err)
			continue
		}
		post, err := db.CreatePost(ctx, database.CreatePostParams{
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: item.Description == ""},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
				continue
			}
			log.Printf("couldn't save post in database: %s", err)
			continue
		}
		log.Printf("saved new post: %s", post.Title)
	}
}
