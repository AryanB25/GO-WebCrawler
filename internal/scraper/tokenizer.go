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

func ExtractWordCounts(data string) map[string]int {
	reader := strings.NewReader(data)
	tokenizer := html.NewTokenizer(reader)
	wordCounts := map[string]int{}

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			break
		}

		if tokenType == html.TextToken {
			tokenData := tokenizer.Token().Data
			tokenFields := strings.Fields(tokenData)
			for _, field := range tokenFields {
				wordCounts[strings.ToLower(field)] = wordCounts[strings.ToLower(field)] + 1
			}
		}
	}

	return wordCounts
}

func ExtractTitle(data string) string {
	reader := strings.NewReader(data)
	tokenizer := html.NewTokenizer(reader)
	var insideTitle bool
	titleString := ""

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()

		if tokenType == html.StartTagToken {
			if token.Data == "title" {
				insideTitle = true
			}
		}

		if tokenType == html.TextToken {
			if insideTitle {
				titleString = token.Data
				insideTitle = false
			}
		}
	}
	return titleString
}
