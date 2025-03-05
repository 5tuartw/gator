package Commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/5tuartw/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func ScrapeFeeds(s *State) error {
	//get the next feed to fetch from db
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get next feed: %w", err)
	}

	//mark as fetched
	if err := s.Db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return fmt.Errorf("failed to mark the feed as fetched: %w", err)
	}

	//fetch the feed using the URL
	rssData, err := FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}
	//iterate items in feed and print title in console
	for _, feedItem := range rssData.Channel.Item {
		published, err := parsePubDate(feedItem.PubDate)
		if err != nil {
			return fmt.Errorf("could not parse publish date of feed: %w", err)
		}
		fmt.Printf("%v: %v\n", rssData.Channel.Title, feedItem.Title)
		err = s.Db.CreatePost(context.Background(), database.CreatePostParams{
			Title:       ConvertNullString(feedItem.Title),
			Url:         feedItem.Link,
			Description: ConvertNullString(feedItem.Description),
			PublishedAt: published,
			FeedID:      feed.ID,
			ID:          uuid.New(),
		})
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" {
					continue
				}
			}
			log.Printf("error inserting post into table: %v", err)
		}
	}

	return nil
}

func ConvertNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func parsePubDate(pubDate string) (time.Time, error) {
	// Try RFC 822 (most common)
	t, err := time.Parse(time.RFC822, pubDate)
	if err == nil {
		return t, nil
	}

	// Try RFC 3339
	t, err = time.Parse(time.RFC3339, pubDate)
	if err == nil {
		return t, nil
	}

	// Try RFC822 with slightly different layouts
	t, err = time.Parse("Mon, 02 Jan 2006 15:04:05 MST", pubDate)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", pubDate)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("02 Jan 2006 15:04:05 MST", pubDate)
	if err == nil {
		return t, nil
	}

	t, err = time.Parse("02 Jan 2006 15:04:05 -0700", pubDate)
	if err == nil {
		return t, nil
	}

	// If all else fails, return an error
	return time.Time{}, fmt.Errorf("failed to parse pubDate: %s", pubDate)
}
