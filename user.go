package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID        string
	Created   int64
	Karma     int
	UpdatedAt string
}

func fetchUser(username string) (*User, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/user/%s.json", username)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("not found")
	}
	var data struct {
		ID      string `json:"id"`
		Created int64  `json:"created"`
		Karma   int    `json:"karma"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &User{
		ID:        data.ID,
		Created:   data.Created,
		Karma:     data.Karma,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func main() {
	csvFile, err := os.Open("bq-results-20260201-205857-1769979839843.csv")
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	usernames := []string{}
	for i, row := range rows {
		if i == 0 || len(row) == 0 {
			continue // skip header
		}
		usernames = append(usernames, strings.TrimSpace(row[0]))
	}

	db, err := sql.Open("sqlite3", "hacker-news.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user (
		id TEXT PRIMARY KEY,
		created INTEGER,
		karma INTEGER,
		updated_at TEXT
	)`)
	if err != nil {
		panic(err)
	}

	// Get all user ids currently in the database
	rowsDb, err := db.Query("SELECT id FROM user")
	if err != nil {
		panic(err)
	}
	defer rowsDb.Close()
	idsInDb := map[string]bool{}
	for rowsDb.Next() {
		var id string
		if err := rowsDb.Scan(&id); err == nil {
			idsInDb[id] = true
		}
	}

	missing := []string{}
	for _, u := range usernames {
		if !idsInDb[u] {
			missing = append(missing, u)
		}
	}
	var processList []string
	if len(missing) > 0 {
		fmt.Printf("Processing %d missing users first...\n", len(missing))
		processList = missing
	} else {
		fmt.Println("All users present in DB, updating all records...")
		processList = usernames
	}

	for _, username := range processList {
		if username == "" {
			fmt.Println("Skipped empty username.")
			continue
		}
		user, err := fetchUser(username)
		if err != nil {
			fmt.Printf("Failed for user %s: %v\n", username, err)
			continue
		}
		// Redacted user check
		if user == nil || user.ID == "" {
			fmt.Printf("Skipped user with empty ID for username: %s\n", username)
			continue
		}
		_, err = db.Exec(`INSERT INTO user (id, created, karma, updated_at) VALUES (?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET created=excluded.created, karma=excluded.karma, updated_at=excluded.updated_at`,
			user.ID, user.Created, user.Karma, user.UpdatedAt)
		if err != nil {
			fmt.Printf("DB error for user %s: %v\n", user.ID, err)
			continue
		}
		fmt.Printf("Upserted user: %s\n", user.ID)
	}
}
