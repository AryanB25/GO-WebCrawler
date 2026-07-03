package scraper

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func FetchData(url string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second} // creates a configurable HTTP connection manager
	connect, err := client.Get(url)                   // HTTP request

	if err != nil { // if there are errors when establishing a HTTP request
		return "", err
	}

	defer connect.Body.Close() // closes the HTTP connection after the function is over to prevent resource leak

	if connect.StatusCode != http.StatusOK { // if the HTTP status is not ok
		return "", fmt.Errorf("Unexpeted Status Code: %v", connect.Status)
	}

	rawData, err := io.ReadAll(connect.Body) // binary data of the url page

	if err != nil { // if there are errors when reading the bytes of the url page
		return "", err
	}

	formattedData := string(rawData) // binary converted to HTML string

	return formattedData, nil // returns the formatted HTML string
}
