package index

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// InitDB opens a connection to the SQLite database and creates the required
// tables if they don't already exist. Returns the database connection for reuse.
func InitDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite", "./gocrawlerdatabase.db") // open or create the database file
	database.SetMaxOpenConns(1)

	if err != nil {
		return database, err
	}

	// create the pages table to store metadata about each crawled page
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS pages (url TEXT PRIMARY KEY, title TEXT, crawled_at TEXT, score INTEGER)")
	if err != nil {
		return database, err
	}

	// create the inverted index table mapping terms to URLs with frequency counts
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS index_table (term TEXT, url TEXT, frequency INTEGER, PRIMARY KEY(term, url))")
	if err != nil {
		return database, err
	}

	return database, err
}

// SavePage inserts a crawled page's metadata into the pages table.
// Records the URL, title, timestamp of when it was crawled, and its keyword score.
func SavePage(db *sql.DB, url string, title string, score int) error {
	_, err := db.Exec("INSERT OR IGNORE INTO pages (url, title, crawled_at, score) VALUES (?, ?, ?, ?)", url, title, time.Now(), score) // insert page with current timestamp
	if err != nil {
		return err
	}
	return nil
}

// SaveTerms inserts word frequency data into the inverted index table.
// For each term in wordCounts, stores a row mapping the term to the page URL
// and how many times that term appeared on that page.
func SaveTerms(db *sql.DB, url string, wordCounts map[string]int) error {
	for term, count := range wordCounts { // iterate over every word found on the page
		_, err := db.Exec("INSERT OR IGNORE INTO index_table (term, url, frequency) VALUES (?, ?, ?)", term, url, count)
		if err != nil {
			return err
		}
	}
	return nil
}

// SearchTerms queries the inverted index for a keyword and returns matching URLs
// ranked by how frequently the term appeared on each page — highest frequency first.
func SearchTerms(db *sql.DB, keyword string) ([]string, error) {
	url_slice := []string{}

	// query the index for all pages containing the keyword, ordered by relevance
	rows, err := db.Query("SELECT url FROM index_table WHERE term = ? ORDER BY frequency DESC", keyword)
	if err != nil {
		return nil, err
	}

	var url string

	for rows.Next() {
		err = rows.Scan(&url) // extract the URL from the current row
		if err != nil {
			return nil, err
		}
		url_slice = append(url_slice, url) // add to results slice
	}

	rows.Close()          // release the database rows resource
	return url_slice, nil // return ranked list of URLs
}
