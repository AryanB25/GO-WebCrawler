package main

import (
	"GO-WebCrawler/internal/datastructures"
	"fmt"
)

func main() {
	sets := datastructures.NewSet()
	sets.Add("https://en.wikipedia.org/wiki/Wikipedia")
	sets.Add("https://aryanbhatt.com")
	fmt.Println(sets.Size())
	fmt.Println(*sets)
	fmt.Println(sets.Contains("https://aryanbhatt.com"))
	fmt.Println(sets.Contains("https://l.com"))
}
