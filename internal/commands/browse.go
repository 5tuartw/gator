package Commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/5tuartw/gator/internal/database"
)

func Browse(s *State, cmd Command, user database.User) error {

	var limit int32 = 2
	if len(cmd.Arguments) > 0 {
		if cmd.Arguments[0] == "--limit" && len(cmd.Arguments) == 2 {
			tryLimit, err := strconv.ParseInt(cmd.Arguments[1], 10, 32)
			if err != nil {
				return fmt.Errorf("second argument must be a valid number")
			}
			limit = int32(tryLimit)
		} else {
			return fmt.Errorf("browse can use the optional argument: --limit <number>")
		}
	}

	posts, err := s.Db.GetUserPosts(context.Background(), database.GetUserPostsParams{
		Name:  s.Config.CurrentUsername,
		Limit: limit,
	})
	if err != nil {
		return fmt.Errorf("unable to fetch posts: %w", err)
	}

	for _, post := range posts {

		description := post.Description.String
		if len(description) > 100 {
			description = description[:100] + "..."
		}

		fmt.Printf("Title:			%v\n", post.Title)
		fmt.Printf("Feed:			%v\n", post.FeedName)
		fmt.Printf("Published:		%v\n", post.PublishedAt)
		fmt.Printf("Description:	%v\n", description)
		fmt.Printf("URL:			%v\n", post.Url)
		fmt.Println("-----------------------------------")
	}

	return nil
}
