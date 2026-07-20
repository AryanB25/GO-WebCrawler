package main

import (
	"GO-WebCrawler/internal/index"
	"GO-WebCrawler/internal/scraper"
	"flag"
	"fmt"
)

func main() {
	searchFlag := flag.String("search", "", "search the crawled index for a keyword")
	crawlTarget := flag.String("target", "books", "keyword to prioritize during crawling")
	flag.Parse()
	if *searchFlag != "" {
		database, err := index.InitDB()
		if err != nil {
			fmt.Println(err)
			return
		}
		listURL, err := index.SearchTerms(database, *searchFlag)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, current := range listURL {
			fmt.Println(current)
		}
	} else {
		scraper.WorkerPool("http://books.toscrape.com", 20, 5, *crawlTarget)
	}
}
