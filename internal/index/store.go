package index

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite", "./gocrawlerdatabase.db")
	if err != nil {
		return database, err
	}
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS pages (url TEXT PRIMARY KEY, title TEXT, crawled_at TEXT, score INTEGER)")
	if err != nil {
		return database, err
	}

	_, err = database.Exec("CREATE TABLE IF NOT EXISTS index_table (term TEXT, url TEXT, frequency INTEGER, PRIMARY KEY(term, url))")
	if err != nil {
		return database, err
	}

	return database, err
}

func SavePage(db *sql.DB, url string, title string, score int) error {
	_, err := db.Exec("INSERT INTO pages (url, title, crawled_at, score) VALUES (?, ?, ?, ?)", url, title, time.Now(), score)
	if err != nil {
		return err
	}
	return nil
}

func SaveTerms(db *sql.DB, url string, wordCounts map[string]int) error {
	for term, count := range wordCounts {
		_, err := db.Exec("INSERT INTO index_table (term, url, frequency) VALUES (?, ?, ?)", term, url, count)
		if err != nil {
			return err
		}
	}
	return nil
}

func SearchTerms(db *sql.DB, keyword string) ([]string, error) {
	url_slice := []string{}

	rows, err := db.Query("SELECT url FROM index_table WHERE term = ? ORDER BY frequency DESC", keyword)
	if err != nil {
		return nil, err
	}

	var url string

	for rows.Next() {
		err = rows.Scan(&url)
		if err != nil {
			return nil, err
		}
		url_slice = append(url_slice, url)
	}

	rows.Close()
	return url_slice, nil
}
