package main

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

func main() {
	text := `<a href="https://a.com">A</a><a href="https://b.com">B</a>` // consists of the HTML
	reader := strings.NewReader(text)
	tokenizer := html.NewTokenizer(reader)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			break
		}

		t := tokenizer.Token()
		if t.Data == "a" {
			if tokenType == html.StartTagToken {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						fmt.Println(attr.Val)
					}
				}
			}
		}
	}
}
