package scraper

import (
	"GO-WebCrawler/internal/datastructures"
	"fmt"
	"net/url"
	"time"
)

func Crawl(seedURL string, maxPages int) {
	startTime := time.Now()
	fmt.Println("Starting crawler...")
	fmt.Println("Maximum Pages:", maxPages)

	queue := datastructures.NewQueue() // create queue
	set := datastructures.NewSet()     // create set
	queue.Enqueue(seedURL)             // enqueue seed URL
	counter := 0                       // keep track of the pages visited

	for !queue.IsEmpty() && counter < maxPages {
		currentURL, notEmpty := queue.Dequeue()

		if !notEmpty { // if the queue is empty
			continue
		}

		if !set.AddIfNotExists(currentURL) {
			continue
		}

		fetchedData, err := FetchData(currentURL) // fetch the HTML data of the page

		if err != nil { // if an error exists
			continue
		}

		listURL := TokenizerURL(fetchedData, 2000) // tokenize the page (string HTML)

		base, err := url.Parse(currentURL)

		if err != nil {
			continue
		}

		for _, pointURL := range listURL {
			relative, err := url.Parse(pointURL)
			if err != nil {
				continue
			}
			resolved := base.ResolveReference(relative)
			queue.Enqueue(resolved.String()) // enqueue the URL's into the queue
		}

		counter++ // increment page count
	}

	elapsed := time.Since(startTime)
	fmt.Println()
	fmt.Printf("Finished crawling %d pages\n", counter)
	fmt.Printf("Runtime: %.2fs\n", elapsed.Seconds())
}
