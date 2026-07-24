package scraper

import (
	"golang.org/x/net/html"
	"strings"
)

// TokenizerURL scans the HTML string and extracts all href values from anchor tags.
// maxTokens limits how many tokens are scanned to prevent slow crawling on large pages.
// Returns a slice of raw URL strings found on the page — may include relative URLs.
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

// ExtractWordCounts scans the HTML string and builds a frequency map of every
// visible word on the page. Only TextTokens are counted — tags and attributes
// are ignored. All words are lowercased for case-insensitive search matching.
// Returns a map of word → number of times it appeared on the page.
func ExtractWordCounts(data string) map[string]int {
	reader := strings.NewReader(data)
	tokenizer := html.NewTokenizer(reader)
	wordCounts := map[string]int{}

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken { // no more tokens to scan
			break
		}

		if tokenType == html.TextToken {
			tokenData := tokenizer.Token().Data      // raw text between HTML tags
			tokenFields := strings.Fields(tokenData) // split into individual words on whitespace
			for _, field := range tokenFields {
				wordCounts[strings.ToLower(field)] = wordCounts[strings.ToLower(field)] + 1 // count each word
			}
		}
	}

	return wordCounts
}

// ExtractTitle scans the HTML string and returns the text content inside the
// first <title> tag found. Uses a boolean flag to detect when the tokenizer
// has entered a title tag and captures the next TextToken as the title.
// Returns an empty string if no title tag is found.
func ExtractTitle(data string) string {
	reader := strings.NewReader(data)
	tokenizer := html.NewTokenizer(reader)
	var insideTitle bool // flag that becomes true once we enter a <title> tag
	titleString := ""

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken { // end of document or parse error
			break
		}

		token := tokenizer.Token()

		if tokenType == html.StartTagToken {
			if token.Data == "title" {
				insideTitle = true // next TextToken will be the title content
			}
		}

		if tokenType == html.TextToken {
			if insideTitle {
				titleString = token.Data // capture the title text
				insideTitle = false      // reset flag — title has been found
			}
		}
	}
	return titleString
}
