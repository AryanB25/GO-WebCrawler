package scraper

import (
	"GO-WebCrawler/internal/datastructures"
	"fmt"
)

func Crawl(seedURL string, maxPages int) {
	queue := datastructures.NewQueue() // create queue
	set := datastructures.NewSet()     // create set
	queue.Enqueue(seedURL)             // enqueue seed URL
	counter := 0                       // keep track of the pages visited

	for !queue.IsEmpty() && counter < maxPages {
		currentURL, notEmpty := queue.Dequeue()

		if !notEmpty { // if the queue is empty
			continue
		}

		if set.Contains(currentURL) { // if the URL has been visited
			continue
		} else {
			set.Add(currentURL) // add the URL to the set if not visited
		}

		fetchedData, err := FetchData(currentURL) // fetch the HTML data of the page

		if err != nil { // if an error exists
			continue
		}

		listURL := TokenizerURL(fetchedData, 2000) // tokenize the page (string HTML)

		fmt.Println("URLs found:", len(listURL))

		for _, url := range listURL {
			queue.Enqueue(url) // enqueue the URL's into the queue
		}

		counter++ // increment page count
	}
}
