# Gator

**This application was made during the Boot.dev course guided project *Build a Blog Aggregator***.


Gator is a command-line application designed to manage user posts and interactions with a PostgreSQL database. It provides various commands to register, login, browse posts, and more.
## Coding skills practiced ##
This application has helped me to practice and improve the following coding skills:
- **Go Programming**: Writing and structuring Go code, using Go modules, and handling errors.
- **Database Management**: Interacting with PostgreSQL, writing SQL queries, and managing database migrations.
- **Command-Line Interfaces**: Building and handling command-line applications, parsing arguments, and providing user feedback.
- **Configuration Management**: Reading and managing configuration files.
- **Concurrency**: Using goroutines and channels for concurrent operations (if applicable).



## Features

- User registration and login
- Browse user posts with optional limit
- Database migrations
- Command-line interface

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/5tuartw/gator.git
   cd gator
2. **Install dependencies:**
Ensure you have Go installed. Then, run:

    ```sh
    go mod tidy
3. **Set up the database:**
Ensure you have PostgreSQL installed and running. Create the database for the project:

    ```sh
    createdb gator
4. **Run database migrations:**
Use the goose tool to run migrations

    ```sh
    goose postgres "postgres://username:password@localhost:5432/gator" up

## Configuration
Create a configuration file config.json in the root directory with the following content:

    {
        DBConnectionString": "postgres://username:password@localhost:5432/gator"
    }
    
Replace username and password with your actual PostrgreSQL username and password

## Usage
Register a new user
```Go
go run . register <username>
```
Login as a user
```Go
go run . login <username>
```
Browse user posts
```
go run . browse [--limit <number>]
```
## Acknowledgements
* Go
* PostgreSQL
* Goose

