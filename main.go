package main

import (
	"fmt"
	"os"

	commands "github.com/alancorleto/blog-aggregator/internal/commands"
	config "github.com/alancorleto/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	state := &commands.State{
		Config: cfg,
	}

	if len(os.Args) < 2 {
		fmt.Println("No command provided.")
		os.Exit(1)
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
