package feeds

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/omn1vor/blog-aggregator/internal/database"
)

func ReadFeeds(ctx context.Context, db *database.Queries, feedList []database.Feed) {
	var wg sync.WaitGroup
	log.Println("Read cycle started")
	for _, feed := range feedList {
		wg.Add(1)
		go func(feed database.Feed) {
			defer wg.Done()

			data, err := Fetch(feed.Url)
			if err != nil {
				log.Printf("Could not fetch rss data from %s: %s\n", feed.Url, err.Error())
				return
			}
			for _, item := range data.Items {
				params := database.CreatePostParams{
					ID:          uuid.New(),
					CreatedAt:   time.Now().UTC(),
					UpdatedAt:   time.Now().UTC(),
					Title:       item.Title,
					Url:         item.Link,
					Description: sql.NullString{String: item.Description, Valid: true},
					PublishedAt: sql.NullTime{Time: parseDate(item.PublishedAt), Valid: true},
					FeedID:      feed.ID,
				}
				_, err := db.CreatePost(ctx, params)
				if err != nil {
					if !strings.Contains(err.Error(), "posts_url_key") {
						log.Println("Could not create a post:", err.Error())
					}
				}
			}
		}(feed)
	}
	wg.Wait()
	log.Println("Read cycle complete")

}

func parseDate(str string) time.Time {
	date := time.Now()
	var err error
	formats := []string{time.UnixDate, time.RFC3339, time.RFC1123, time.RFC1123Z}
	for _, format := range formats {
		date, err = time.Parse(format, str)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Println("Error while parsing date:", str)
	}
	return date
}
