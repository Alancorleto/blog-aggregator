package commands

import (
	"context"
	"fmt"
	"time"

	database "github.com/alancorleto/blog-aggregator/internal/database"
	feedfetcher "github.com/alancorleto/blog-aggregator/internal/feed_fetcher"
	state "github.com/alancorleto/blog-aggregator/internal/state"
	"github.com/google/uuid"
)

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	CommandsMap map[string]func(*state.State, Command) error
}

func InitializeCommands() *Commands {
	cmds := &Commands{
		CommandsMap: make(map[string]func(*state.State, Command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)

	return cmds
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	if handler, exists := c.CommandsMap[cmd.Name]; exists {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.Name)
}

func (c *Commands) register(name string, handler func(*state.State, Command) error) {
	c.CommandsMap[name] = handler
}

func handlerLogin(state *state.State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("username argument is required for login command")
	}

	userName := cmd.Arguments[0]

	if _, err := state.Db.GetUser(context.Background(), userName); err != nil {
		return fmt.Errorf("user '%s' does not exist", userName)
	}

	err := state.Config.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' logged in successfully.\n", userName)
	return nil
}

func handlerRegister(state *state.State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("username argument is required for register command")
	}

	userName := cmd.Arguments[0]

	if _, err := state.Db.GetUser(context.Background(), userName); err == nil {
		return fmt.Errorf("user '%s' already exists", userName)
	}

	state.Db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			Name:      userName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	)

	err := state.Config.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' registered successfully.\n", userName)
	return nil
}

func handlerReset(state *state.State, cmd Command) error {
	err := state.Db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset users: %v", err)
	}

	fmt.Println("All users have been reset successfully.")
	return nil
}

func handlerUsers(state *state.State, cmd Command) error {
	users, err := state.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}

	loggedUserName := state.Config.CurrentUserName

	fmt.Println("Registered users:")
	for _, user := range users {
		if user == loggedUserName {
			fmt.Println("*", user, "(current)")
		} else {
			fmt.Println("*", user)
		}
	}

	return nil
}

func handlerAgg(state *state.State, cmd Command) error {
	// if len(cmd.Arguments) < 1 {
	// 	return fmt.Errorf("feed URL argument is required for fetchfeed command")
	// }

	feedURL := "https://www.wagslane.dev/index.xml"

	feed, err := feedfetcher.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %v", err)
	}

	fmt.Printf("%+v\n", feed)

	return nil
}
