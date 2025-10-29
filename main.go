package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	commands "github.com/alancorleto/blog-aggregator/internal/commands"
	config "github.com/alancorleto/blog-aggregator/internal/config"
	database "github.com/alancorleto/blog-aggregator/internal/database"
	state "github.com/alancorleto/blog-aggregator/internal/state"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("No command provided.")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	defer db.Close()
	dbQueries := database.New(db)

	state := &state.State{
		Config: cfg,
		Db:     dbQueries,
	}

	cmd := commands.Command{
		Name:      os.Args[1],
		Arguments: os.Args[2:],
	}

	cmds := commands.InitializeCommands()
	err = cmds.Run(state, cmd)
	if err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}
