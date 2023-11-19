package main

import (
	"context"
	"log"
	"time"

	"github.com/omn1vor/blog-aggregator/internal/feeds"
)

func StartReadingFeeds(ctx context.Context, cfg apiConfig) {
	ticker := time.NewTicker(cfg.ReadingInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("Received ticker signal")
				readFeeds(ctx, cfg)
			case <-ctx.Done():
				log.Println("Received cancel signal")
				ticker.Stop()
			}
		}
	}()
}

func readFeeds(ctx context.Context, cfg apiConfig) {
	feedList, err := cfg.DB.GetNextFeedsToFetch(context.Background(), int32(cfg.NumberOfFeedsToRead))
	if err != nil {
		log.Println("Error while getting next feeds to fetch", err.Error())
		return
	}
	feeds.ReadFeeds(ctx, cfg.DB, feedList)
}
