package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/mxrb/gator/internal/config"
	"github.com/mxrb/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	// log.Println("Starting Gopher ...")
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}
	db, err := sql.Open("postgres", cfg.DbUrl)
	dbQueries := database.New(db)

	programState := state{cfg: &cfg, db: dbQueries}
	cmds := commands{
		registeredCommands: make(map[string]commandHandler),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddfeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		fmt.Println("Usage: gator <command> [args...]")
	}
	cmd := command{Name: os.Args[1], Args: os.Args[2:]}
	if err := cmds.run(&programState, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type state struct {
	cfg *config.Config
	db  *database.Queries
}
