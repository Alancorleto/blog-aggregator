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
	cmds.register("addfeed", middleWareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middleWareLoggedIn(handlerFollow))
	cmds.register("following", middleWareLoggedIn(handlerFollowing))

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

func middleWareLoggedIn(handler func(state *state.State, cmd Command, user database.User) error) func(state *state.State, cmd Command) error {
	return func(state *state.State, cmd Command) error {
		currentUserName := state.Config.CurrentUserName
		user, err := state.Db.GetUser(context.Background(), currentUserName)
		if err != nil {
			return err
		}
		return handler(state, cmd, user)
	}
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

	err = state.Db.ResetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset feeds: %v", err)
	}

	err = state.Db.ResetFeedFollows(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset feed_follows: %v", err)
	}

	fmt.Println("All databases have been reset successfully.")
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

func handlerAddFeed(state *state.State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("not enough arguments for add feed command, expected 2, got %d", len(cmd.Arguments))
	}

	feedName := cmd.Arguments[0]
	feedUrl := cmd.Arguments[1]

	feed, err := state.Db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			Name:      feedName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Url:       feedUrl,
			UserID:    user.ID,
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", feed)

	followFeed(user, feed.Url, state.Db)

	return nil
}

func handlerFeeds(state *state.State, cmd Command) error {
	feeds, err := state.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("--- %s ---\nURL: %s\nUser: %s\n\n", feed.Name, feed.Url, feed.UserName)
	}

	return nil
}

func handlerFollow(state *state.State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("expected 1 argument, got 0")
	}

	feedUrl := cmd.Arguments[0]

	feedFollowResponse, err := followFeed(user, feedUrl, state.Db)
	if err != nil {
		return err
	}

	fmt.Printf("%s is now following %s\n", feedFollowResponse.UserName, feedFollowResponse.FeedName)

	return nil
}

func followFeed(user database.User, feedUrl string, db *database.Queries) (database.CreateFeedFollowRow, error) {
	feed, err := db.GetFeedByURL(context.Background(), feedUrl)
	if err != nil {
		return database.CreateFeedFollowRow{}, err
	}

	feedFollowResponse, err := db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return database.CreateFeedFollowRow{}, err
	}

	return feedFollowResponse, nil
}

func handlerFollowing(state *state.State, cmd Command, user database.User) error {
	feedFollows, err := state.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("Feeds followed by %s:\n", user.Name)

	for _, feedFollow := range feedFollows {
		fmt.Printf("- %s\n", feedFollow.FeedName)
	}

	return nil
}
