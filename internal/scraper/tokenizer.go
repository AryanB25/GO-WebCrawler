package scraper

import (
	"golang.org/x/net/html"
	"strings"
)

func TokenizerURL(data string, maxTokens int) []string {
	reader := strings.NewReader(data)
	tokenizer := html.NewTokenizer(reader)
	listUrl := []string{} // initializing the slice containing the URL's
	tokenCount := 0       // counter for the number of tokens

	for tokenCount <= maxTokens { // run until the token count is less than the allowed limit

		tokenType := tokenizer.Next() // find the token type of the data

		if tokenType == html.ErrorToken { // if there is no more data or HTML pieces
			break
		}

		tokenCount++ // increment the token counter

		token := tokenizer.Token() // contains the current token

		if token.Data == "a" {
			if tokenType == html.StartTagToken {
				for _, attr := range token.Attr {
					if attr.Key == "href" { // if the key is an "href" for links
						listUrl = append(listUrl, attr.Val) // add the URL to the slice
					}
				}
			}
		}
	}

	return listUrl // returns the links found in the HTML pieces or data
}
