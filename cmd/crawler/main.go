package main

import "GO-WebCrawler/internal/scraper"

func main() {
	scraper.Crawl("http://books.toscrape.com", 50)
}
