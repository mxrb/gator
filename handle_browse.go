package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"text/template"

	"github.com/mxrb/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	switch len(cmd.Args) {
	case 0:
		limit = 2
	case 1:
		var err error
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("argument cannot be converted to integer: %w", err)
		}
	default:
		return fmt.Errorf("usage: %s <limit>", cmd.Name)
	}
	posts, err := s.db.ListPostsForUser(context.Background(), database.ListPostsForUserParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		return fmt.Errorf("couldn't fetch followed feeds: %w", err)
	}
	const templText = "{{range .}}* {{.Title}}\n{{else}}No posts found!\n{{end}}"
	t := template.Must(template.New("following").Parse(templText))
	if err := t.Execute(os.Stdout, posts); err != nil {
		return fmt.Errorf("couldn't render output: %w", err)
	}
	return nil
}
