package main

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/mxrb/gator/internal/database"
)

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.ListFeedsWithUserName(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list feeds: %w", err)
	}
	printFeedsResult(feeds)
	return nil
}

func printFeedsResult(feeds []database.ListFeedsWithUserNameRow) error {
	const templText = `{{range .}}Name: {{.Name}}
Url:  {{.Url}}
User: {{.UserName}}
===============================================
{{else}}
There are no feeds in the database.
{{end}}
`
	t := template.Must(template.New("feeds").Parse(templText))
	return t.Execute(os.Stdout, feeds)
}
