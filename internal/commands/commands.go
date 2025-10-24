package commands

import (
	"fmt"

	config "github.com/alancorleto/blog-aggregator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	CommandsMap map[string]func(*State, Command) error
}

func InitializeCommands() *Commands {
	cmds := &Commands{
		CommandsMap: make(map[string]func(*State, Command) error),
	}

	cmds.register("login", handlerLogin)

	return cmds
}

func (c *Commands) Run(s *State, cmd Command) error {
	if handler, exists := c.CommandsMap[cmd.Name]; exists {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.Name)
}

func (c *Commands) register(name string, handler func(*State, Command) error) {
	c.CommandsMap[name] = handler
}

func handlerLogin(state *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("username argument is required for login command")
	}

	userName := cmd.Arguments[0]
	err := state.Config.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' logged in successfully.\n", userName)
	return nil
}
