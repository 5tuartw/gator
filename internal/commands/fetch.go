package Commands

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/5tuartw/gator/internal/database"
)

func FetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {

	newReq, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err)
		return nil, err
	}

	newReq.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: non-200 status code: ", resp.StatusCode)
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	//read the response body with io.ReadAll
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return nil, err
	}

	//unmarshal the xml data into an RSSFeed struct
	var xmlData RSSFeed
	err = xml.Unmarshal(data, &xmlData)
	if err != nil {
		fmt.Println("Error marshalling xml data: ", err)
		return nil, err
	}

	//unescape HTML entities in the title and description fields
	xmlData.Channel.Title = html.UnescapeString(xmlData.Channel.Title)
	xmlData.Channel.Description = html.UnescapeString(xmlData.Channel.Description)
	for i := range xmlData.Channel.Item {
		xmlData.Channel.Item[i].Title = html.UnescapeString(xmlData.Channel.Item[i].Title)
		xmlData.Channel.Item[i].Description = html.UnescapeString(xmlData.Channel.Item[i].Description)
	}

	return &xmlData, nil
}

func AggCommand(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("agg requires 1 argument: agg <time>")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("failed to parse time argument: %w", err)
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		ScrapeFeeds(s)
	}
}

// ESCAPED CODE FROM ORIGINAL AGG
// Create context
//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//defer cancel()

// Define the feed URL
//feedURL := "https://www.wagslane.dev/index.xml"

// Call your existing FetchFeed function
//feed, err := FetchFeed(ctx, feedURL)
//if err != nil {
//	return err
//}

// Print the entire struct
//fmt.Printf("%+v\n", feed)
//return nil

func AddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		return fmt.Errorf("adding a feed requires 2 arguments <name> and <url>")
	}

	fName := cmd.Arguments[0]
	fURL := cmd.Arguments[1]
	//fmt.Printf("current user %v\n", userName)

	userId := user.ID

	newFeed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:   fName,
		Url:    fURL,
		UserID: userId,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: userId,
		FeedID: newFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed follow entry")
	}

	fmt.Println("Feed created successfully")
	printFeed(newFeed)

	return nil
}

func Feeds(s *State, cmd Command) error {
	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("feeds command takes no arguments")
	}

	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %w", err)
	}
	for _, feed := range feeds {
		username, err := s.Db.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error looking up user id")
		}
		fmt.Println("=========================")
		fmt.Printf("Name:            %s\n", feed.Name)
		fmt.Printf("URL:             %s\n", feed.Url)
		fmt.Printf("User:            %s\n", username)
		fmt.Println("=========================")
	}
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Println("=========================")
	fmt.Printf("Name:            %s\n", feed.Name)
	fmt.Printf("URL:             %s\n", feed.Url)
	fmt.Printf("User:            %s\n", feed.UserID)
	fmt.Println("=========================")
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}
