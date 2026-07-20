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

func WorkerPool(seedURL string, limit int, numberWorkers int, target string) {
	queue := datastructures.NewPriorityQueue()
	set := datastructures.NewSet()
	score := strings.Count(seedURL, target)
	queue.Push(seedURL, score)
	database, err := index.InitDB()

	if err != nil {
		panic(err)
	}

	defer database.Close()

	jobs := make(chan string, numberWorkers*2)
	var wg sync.WaitGroup
	var mutex sync.Mutex

	counter := 0

	// worker goroutine
	for range numberWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for currentURL := range jobs {
				if set.Contains(currentURL) { // if the URL has been visited
					continue
				} else {
					set.Add(currentURL) // add the URL to the set if not visited
				}

				fetchedData, err := FetchData(currentURL) // fetch the HTML data of the page

				if err != nil { // if an error exists
					continue
				}

				titleURL := ExtractTitle(fetchedData)
				wordCount := ExtractWordCounts(fetchedData)
				pageScore := strings.Count(currentURL, target)

				err = index.SavePage(database, currentURL, titleURL, pageScore)

				if err != nil {
					continue
				}

				err = index.SaveTerms(database, currentURL, wordCount)

				if err != nil {
					continue
				}

				listURL := TokenizerURL(fetchedData, 2000) // tokenize the page (string HTML)

				fmt.Println("URLs found:", len(listURL))

				base, err := url.Parse(currentURL)

				if err != nil {
					continue
				}

				for _, foundURL := range listURL {
					relative, err := url.Parse(foundURL)

					if err != nil {
						continue
					}

					resolved := base.ResolveReference(relative)
					urlScore := strings.Count(resolved.String(), target)
					queue.Push(resolved.String(), urlScore)
				}

				mutex.Lock()
				counter++ // increment page count
				mutex.Unlock()
			}
		}()
	}

	// feeder goroutine
	go func() {
		for {
			mutex.Lock()
			c := counter
			mutex.Unlock()

			if c > limit {
				break
			}

			if queue.IsEmpty() {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			currentURL, notEmpty := queue.Pop() // extract the first priority item

			if !notEmpty { // if the queue is empty
				continue
			} else {
				jobs <- currentURL
			}
		}
		close(jobs)
	}()

	wg.Wait() // waits for all the workers to finish
}
