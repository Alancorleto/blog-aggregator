package feedscraper

import (
	"context"
	"database/sql"
	"strings"
	"time"

	database "github.com/alancorleto/blog-aggregator/internal/database"
	feedfetcher "github.com/alancorleto/blog-aggregator/internal/feed_fetcher"
	"github.com/google/uuid"
)

func ScrapeNextFeed(db *database.Queries) (string, error) {
	nextFeed, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return "", err
	}

	err = db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return "", err
	}

	rssFeed, err := feedfetcher.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return "", err
	}

	for _, rssItem := range rssFeed.Channel.Item {
		rssItemPubDate, err := time.Parse(time.RFC1123, rssItem.PubDate)
		if err != nil {
			return "", err
		}
		_, err = db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       rssItem.Title,
				Url:         rssItem.Link,
				Description: sql.NullString{String: rssItem.Description, Valid: rssItem.Description != ""},
				PublishedAt: rssItemPubDate,
				FeedID:      nextFeed.ID,
			},
		)
		if err != nil && !strings.Contains(err.Error(), "posts_url_key") {
			return "", err
		}
	}

	return rssFeed.Channel.Title, nil
}
