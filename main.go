package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/j-machuca/gator/internal/config"
	"github.com/j-machuca/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	programState := &state{
		cfg: &cfg,
	}
	programState.db = dbQueries

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handleLogin)
	cmds.register("register", handleRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleUsers)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFeedFollow))
	cmds.register("following", middlewareLoggedIn(handleGetFeedsForUser))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))

	cmds.register("agg", handleAgg)
	if len(os.Args) < 2 {
		log.Fatal("usage cli <command> [args...]")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
