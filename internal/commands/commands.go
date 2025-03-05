package Commands

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	Config "github.com/5tuartw/gator/internal/config"
	"github.com/5tuartw/gator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	Db     *database.Queries
	Config *Config.Config
}

type Command struct {
	Name      string
	Arguments []string
}

// function to handle commands which require a user (current logged in user)
func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUsername)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

// function to log in a new user from the users table
func HandlerLogin(s *State, c Command) error {
	if len(c.Arguments) != 1 {
		return fmt.Errorf("login requires exactly one argument")
	}
	username := c.Arguments[0]

	//check if user exists
	_, err := s.Db.GetUser(context.Background(), username)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			fmt.Printf("User with name '%s' does not exist\n", username)
			os.Exit(1)
		}
		return err
	}

	err = s.Config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("User has been set to: %s\n", username)
	return nil
}

// function to register a new user in the users table
func RegisterHandler(s *State, c Command) error {
	if len(c.Arguments) < 1 {
		return fmt.Errorf("register requires a name argument")
	}

	name := c.Arguments[0]
	user, err := s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	// Handle errors - particularly if user already exists
	if err != nil {
		// Check if the error is because the user already exists
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			fmt.Printf("User with name '%s' already exists\n", name)
			os.Exit(1)
		}
		// For other errors, return them
		return err
	}

	s.Config.SetUser(name)

	fmt.Printf("User has been registered: %s\n", user.Name)
	//fmt.Printf("User details: %+v\n", user)

	return nil
}

// function to reset the users database
func Reset(s *State, c Command) error {
	if len(c.Arguments) > 0 {
		return fmt.Errorf("reset does not take any arguments")
	}
	// Reset the database
	err := s.Db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset user table: %w", err)
	}

	fmt.Println("Users table has been reset")
	return nil
}

// function to list the current users
func Users(s *State, c Command) error {
	if len(c.Arguments) > 0 {
		return fmt.Errorf("users does not take any arguments")
	}

	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	tag := "(current)"

	fmt.Println("Users:")
	for _, user := range users {
		if user.Name == s.Config.CurrentUsername {
			fmt.Printf("* %s %s\n", user.Name, tag)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}
	return nil
}

type Commands struct {
	Commands map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Commands[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	if f, ok := c.Commands[cmd.Name]; ok {
		return f(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.Name)
}
