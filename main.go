package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	Commands "github.com/5tuartw/gator/internal/commands"

	Config "github.com/5tuartw/gator/internal/config"
	"github.com/5tuartw/gator/internal/database"
)

//psql "postgres://postgres:postgres@localhost:5432/gator"
//goose postgres "postgres://postgres:postgres@localhost:5432/gator" up

func main() {
	cfg, err := Config.Read()
	if err != nil {
		fmt.Println("Error reading config: ", err)
		return
	}

	//fmt.Println("Config read successfully: ", cfg)
	//setting up state
	state := &Commands.State{Config: &cfg}

	//setting up database connection
	db, err := sql.Open("postgres", cfg.DBConnectionString)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Run migrations before using the database
	if err := RunMigration(db); err != nil {
		fmt.Println("Failed to run migrations:", err)
		os.Exit(1)
	}
	// Check if the database connection is valid
	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	state.Db = dbQueries

	cmds := &Commands.Commands{Commands: make(map[string]func(*Commands.State, Commands.Command) error)}

	cmds.Register("login", Commands.HandlerLogin)
	cmds.Register("register", Commands.RegisterHandler)
	cmds.Register("reset", Commands.Reset)
	cmds.Register("users", Commands.Users)
	cmds.Register("agg", Commands.AggCommand)
	cmds.Register("feeds", Commands.Feeds)

	//commands requiring logging in
	cmds.Register("addfeed", Commands.MiddlewareLoggedIn(Commands.AddFeed))
	cmds.Register("follow", Commands.MiddlewareLoggedIn(Commands.Follow))
	cmds.Register("following", Commands.MiddlewareLoggedIn(Commands.Following))
	cmds.Register("unfollow", Commands.MiddlewareLoggedIn(Commands.Unfollow))
	cmds.Register("browse", Commands.MiddlewareLoggedIn(Commands.Browse))

	//cmd := Commands.Command{Name: "login", Arguments: []string{"bob"}}
	if len(os.Args) < 2 {
		fmt.Println("Too few arguments")
		os.Exit(1)
	}
	if state.Db == nil {
		fmt.Println("Database connection is not initialized")
		os.Exit(1)
	}
	cmd := Commands.Command{Name: os.Args[1], Arguments: os.Args[2:]}
	err = cmds.Run(state, cmd)
	if err != nil {
		fmt.Println("Error running command: ", err)
		os.Exit(1)
	}

	//fmt.Println("Command executed successfully")

}
