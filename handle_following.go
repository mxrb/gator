package main

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/mxrb/gator/internal/database"
)

func handlerFollowing(s *state, _ command, user database.User) error {
	feeds, err := s.db.ListFeedsFollowedByUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't fetch followed feeds: %w", err)
	}
	const templText = "{{range .}}* {{.Name}}\n{{else}}You are not following any feeds.{{end}}"
	t := template.Must(template.New("following").Parse(templText))
	if err := t.Execute(os.Stdout, feeds); err != nil {
		return fmt.Errorf("couldn't render output: %w", err)
	}
	return nil
}
