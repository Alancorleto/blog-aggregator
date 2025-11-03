package feedscraper

import (
	"context"
	"fmt"

	database "github.com/alancorleto/blog-aggregator/internal/database"
	feedfetcher "github.com/alancorleto/blog-aggregator/internal/feed_fetcher"
)

func ScrapeFeeds(db *database.Queries) error {
	nextFeed, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return err
	}

	rssFeed, err := feedfetcher.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	fmt.Printf("-- Printing titles from feed: %s --\n", rssFeed.Channel.Title)
	for _, rssItem := range rssFeed.Channel.Item {
		fmt.Println(rssItem.Title)
	}

	return nil
}
