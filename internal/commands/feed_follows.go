package Commands

import (
	"context"
	"fmt"

	"github.com/5tuartw/gator/internal/database"
)

func Follow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("follow requires only one argument: follow <url>")
	}
	fURL := cmd.Arguments[0]

	userId := user.ID

	thisFeed, err := s.Db.GetFeedId(context.Background(), fURL)
	if err != nil {
		return fmt.Errorf("feed does not exist in feeds table %w", err)
	}
	feedId := thisFeed.ID

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: userId,
		FeedID: feedId,
	})
	if err != nil {
		return fmt.Errorf("unable to create new feed follow %w", err)
	}
	fmt.Printf("New feed follow (id: %v) for user: %s feed: %s\n", feedFollow.ID, user.Name, fURL)

	return nil
}

func Unfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("unfollow requires only one argument: unfollow <url>")
	}

	err := s.Db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		Url:    cmd.Arguments[0],
	})
	if err != nil {
		return fmt.Errorf("could not unfollow feed")
	}

	return nil
}

func Following(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("following command requires no arguments")
	}

	//fmt.Printf("current user %v\n", userName)

	followingFeeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting follower feeds: %w", err)
	}

	fmt.Printf("User: %s is following:\n", user.Name)
	for _, feed := range followingFeeds {
		thisFeed, err := s.Db.GetFeedName(context.Background(), feed.FeedID)
		if err != nil {
			return fmt.Errorf("error finding feed name in table")
		}
		fmt.Printf("        - %v\n", thisFeed)
	}

	return nil
}
