package main

import "GO-WebCrawler/internal/scraper"

func main() {
	scraper.WorkerPool("http://books.toscrape.com", 20, 5, "aryan")
}
