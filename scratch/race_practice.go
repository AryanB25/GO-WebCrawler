package main

import (
	"GO-WebCrawler/internal/datastructures"
	"sync"
)

func main() {
	queue := datastructures.NewQueue()
	var wg sync.WaitGroup

	for j := 1; j <= 5; j++ {
		wg.Add(1)
		go func(queue *datastructures.Queue) {
			defer wg.Done()
			for i := 1; i <= 100; i++ {
				queue.Enqueue("https://example.com")
			}
		}(queue)
	}

	wg.Wait()
}
