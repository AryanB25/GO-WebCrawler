package scraper

import (
	"GO-WebCrawler/internal/datastructures"
	"GO-WebCrawler/internal/index"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"
)

// WorkerPool runs a concurrent web crawler using a fixed pool of worker goroutines.
// It fetches pages starting from seedURL, prioritizing pages that contain the target keyword,
// and persists all crawled data to a SQLite database for later searching.
// seedURL is the starting page, limit is the max number of pages to crawl,
// numberWorkers controls how many goroutines fetch pages simultaneously,
// and target is the keyword used to score and prioritize URLs.
func WorkerPool(seedURL string, limit int, numberWorkers int, target string) {
	queue := datastructures.NewPriorityQueue() // priority queue orders URLs by keyword relevance score
	set := datastructures.NewSet()             // set tracks visited URLs to prevent revisiting
	score := strings.Count(seedURL, target)    // score the seed URL before pushing
	queue.Push(seedURL, score)

	// initialize the database connection and create tables if they don't exist
	database, err := index.InitDB()
	if err != nil {
		panic(err)
	}
	defer database.Close() // close the database connection when WorkerPool exits

	jobs := make(chan string, numberWorkers*2) // buffered channel distributes URLs to workers
	var wg sync.WaitGroup
	var mutex sync.Mutex

	counter := 0 // tracks how many pages have been successfully crawled

	// launch N worker goroutines — each independently fetches and processes pages
	for range numberWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done() // signal this worker is done when it exits
			for currentURL := range jobs {
				if set.Contains(currentURL) { // skip already visited URLs
					continue
				} else {
					set.Add(currentURL) // mark URL as visited before fetching
				}

				fetchedData, err := FetchData(currentURL) // fetch raw HTML from the page
				if err != nil {                           // skip this URL if fetch fails
					continue
				}

				titleURL := ExtractTitle(fetchedData)          // extract the page title from HTML
				wordCount := ExtractWordCounts(fetchedData)    // build word frequency map for the inverted index
				pageScore := strings.Count(currentURL, target) // score based on keyword in URL

				// persist page metadata to the pages table
				err = index.SavePage(database, currentURL, titleURL, pageScore)
				if err != nil {
					continue
				}

				// persist word frequencies to the inverted index table
				err = index.SaveTerms(database, currentURL, wordCount)
				if err != nil {
					continue
				}

				listURL := TokenizerURL(fetchedData, 2000) // extract all hrefs from the page HTML
				fmt.Println("URLs found:", len(listURL))

				base, err := url.Parse(currentURL) // parse base URL for resolving relative links
				if err != nil {
					continue
				}

				for _, foundURL := range listURL {
					relative, err := url.Parse(foundURL) // parse each found href
					if err != nil {
						continue
					}

					resolved := base.ResolveReference(relative)          // turn relative URLs into absolute ones
					urlScore := strings.Count(resolved.String(), target) // score the resolved URL by keyword
					queue.Push(resolved.String(), urlScore)              // push to priority queue with score
				}

				mutex.Lock()
				counter++ // safely increment page count across goroutines
				mutex.Unlock()
			}
		}()
	}

	// feeder goroutine — moves URLs from the priority queue into the jobs channel
	// runs concurrently with workers so crawling continues as new URLs are discovered
	go func() {
		for {
			mutex.Lock()
			c := counter // safely read counter
			mutex.Unlock()

			if c > limit { // stop feeding once page limit is reached
				break
			}

			if queue.IsEmpty() {
				time.Sleep(100 * time.Millisecond) // wait briefly if queue is temporarily empty
				continue
			}

			currentURL, notEmpty := queue.Pop() // get highest priority URL from the heap
			if !notEmpty {                      // safety check in case queue emptied between IsEmpty and Pop
				continue
			} else {
				jobs <- currentURL // send URL to next available worker
			}
		}
		close(jobs) // closing the channel signals all workers to stop ranging and exit
	}()

	wg.Wait() // block until all workers have finished processing
}
