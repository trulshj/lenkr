package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func ConnectToDb() *sql.DB {
	url := os.Getenv("DB_URL")

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	return db
}

func GetLink(db *sql.DB, shortUrl string) (Link, error) {
	var link Link

	row := db.QueryRow("SELECT * FROM links WHERE short_url = ?", shortUrl)
	if err := row.Scan(&link.Url, &link.ShortUrl, &link.Fetches); err != nil {
		return link, fmt.Errorf("failed to get link: %v", err)
	}

	return link, nil
}

func GetLinks(db *sql.DB) []Link {
	rows, err := db.Query("SELECT * FROM links")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to query db: %v\n", err)
		os.Exit(1)
	}

	defer rows.Close()

	var links []Link

	for rows.Next() {
		var link Link
		if err := rows.Scan(&link.Url, &link.ShortUrl, &link.Fetches); err != nil {
			fmt.Println("Error scanning row: ", err)
		}

		links = append(links, link)
		fmt.Printf("Link: %+v\n", link)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during row iteration: ", err)
	}

	fmt.Println("Links: ", links)

	return links
}

func IncrementFetches(db *sql.DB, link Link) error {
	_, err := db.Exec("UPDATE links SET fetches = ? WHERE short_url = ?", link.Fetches+1, link.ShortUrl)
	if err != nil {
		return fmt.Errorf("failed to increment fetches: %v", err)
	}

	return nil
}

func InsertLink(db *sql.DB, url string, shortUrl string) error {
	if url == "" || shortUrl == "" {
		return errors.New("url and shortUrl cannot be empty")
	}

	_, err := db.Exec("INSERT INTO links (url, short_url) VALUES (?, ?)", url, shortUrl)
	if err != nil {
		return fmt.Errorf("failed to insert link: %v", err)
	}

	return nil
}
