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
	startTime := time.Now()
	fmt.Println("Starting crawler...")
	fmt.Println("Workers:", numberWorkers)
	fmt.Println("Limit:", limit)

	queue := datastructures.NewPriorityQueue() // priority queue orders URLs by keyword relevance score
	set := datastructures.NewSet()             // set tracks visited URLs to prevent revisiting
	score := strings.Count(seedURL, target)    // score the seed URL before pushing
	queue.Push(seedURL, score)

	var mutex sync.Mutex

	// initialize the database connection and create tables if they don't exist
	database, err := index.InitDB()
	if err != nil {
		panic(err)
	}
	defer database.Close() // close the database connection when WorkerPool exits

	jobs := make(chan string, numberWorkers) // buffered channel distributes URLs to workers
	var wg sync.WaitGroup

	counter := 0 // tracks how many pages have been successfully crawled

	// launch N worker goroutines — each independently fetches and processes pages
	for range numberWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done() // signal this worker is done when it exits
			for currentURL := range jobs {
				mutex.Lock()
				if counter >= limit {
					mutex.Unlock()
					continue
				}
				mutex.Unlock()

				if !set.AddIfNotExists(currentURL) {
					continue
				}

				fetchedData, err := FetchData(currentURL) // fetch raw HTML from the page
				if err != nil {
					fmt.Println("fetch error:", err)
					continue
				}

				titleURL := ExtractTitle(fetchedData)          // extract the page title from HTML
				wordCount := ExtractWordCounts(fetchedData)    // build word frequency map for the inverted index
				pageScore := strings.Count(currentURL, target) // score based on keyword in URL

				mutex.Lock()
				// persist page metadata to the pages table
				err = index.SavePage(database, currentURL, titleURL, pageScore)
				if err == nil {
					// persist word frequencies to the inverted index table
					err = index.SaveTerms(database, currentURL, wordCount)
				}

				mutex.Unlock()

				if err != nil {
					continue
				}

				listURL := TokenizerURL(fetchedData, 2000) // extract all hrefs from the page HTML

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
			currentURL, notEmpty := queue.Pop()

			if !notEmpty {
				continue
			}

			mutex.Lock()

			if counter >= limit {
				mutex.Unlock()
				break
			}

			counter++
			mutex.Unlock()
			jobs <- currentURL
		}
		close(jobs) // closing the channel signals all workers to stop ranging and exit
	}()

	wg.Wait() // block until all workers have finished processing
	elapsed := time.Since(startTime)
	fmt.Println()
	fmt.Printf("Finished crawling %d pages\n", counter)
	fmt.Printf("Runtime: %.2fs\n", elapsed.Seconds())
}
