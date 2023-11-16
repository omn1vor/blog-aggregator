package feeds

import (
	"log"
	"sync"
)

func ReadFeeds(urls []string) {
	var wg sync.WaitGroup
	log.Println("Read cycle started")
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			data, err := Fetch(url)
			if err != nil {
				log.Printf("Could not fetch rss data from %s: %s\n", url, err.Error())
				return
			}
			for _, item := range data.Items {
				log.Println(item.Title, ":", item.Link)
			}
		}(url)
	}
	wg.Wait()
	log.Println("Read cycle complete")

}
