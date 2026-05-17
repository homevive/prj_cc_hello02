package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

// GitHubItem 统一存储 Issue 和 PR
type GitHubItem struct {
	ID        int64     `json:"id"`
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	ItemType  string    `json:"item_type"`
	Author    string    `json:"author"`
	Labels    string    `json:"labels"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Stats 统计信息
type Stats struct {
	IssuesOpen   int `json:"issues_open"`
	IssuesClosed int `json:"issues_closed"`
	PRsOpen      int `json:"prs_open"`
	PRsClosed    int `json:"prs_closed"`
	PRsMerged    int `json:"prs_merged"`
	TotalItems   int `json:"total_items"`
}

// Snapshot 历史快照
type Snapshot struct {
	ID          int       `json:"id"`
	SnapshotAt  time.Time `json:"snapshot_at"`
	IssuesOpen  int       `json:"issues_open"`
	IssuesClosed int      `json:"issues_closed"`
	PRsOpen     int       `json:"prs_open"`
	PRsClosed   int       `json:"prs_closed"`
	PRsMerged   int       `json:"prs_merged"`
	TotalItems  int       `json:"total_items"`
}

var db *sql.DB

func initDB(path string) error {
	var err error
	db, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS github_items (
			id INTEGER PRIMARY KEY,
			number INTEGER NOT NULL,
			title TEXT NOT NULL,
			state TEXT NOT NULL,
			item_type TEXT NOT NULL,
			author TEXT NOT NULL DEFAULT '',
			labels TEXT NOT NULL DEFAULT '',
			url TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			UNIQUE(number, item_type)
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_item_type ON github_items(item_type)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_state ON github_items(state)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS snapshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			snapshot_at DATETIME NOT NULL,
			issues_open INTEGER NOT NULL DEFAULT 0,
			issues_closed INTEGER NOT NULL DEFAULT 0,
			prs_open INTEGER NOT NULL DEFAULT 0,
			prs_closed INTEGER NOT NULL DEFAULT 0,
			prs_merged INTEGER NOT NULL DEFAULT 0,
			total_items INTEGER NOT NULL DEFAULT 0
		)
	`)
	return err
}

func upsertItem(item GitHubItem) error {
	_, err := db.Exec(`
		INSERT INTO github_items (id, number, title, state, item_type, author, labels, url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(number, item_type) DO UPDATE SET
			state=excluded.state, title=excluded.title, labels=excluded.labels,
			updated_at=excluded.updated_at
	`, item.ID, item.Number, item.Title, item.State, item.ItemType, item.Author, item.Labels, item.URL, item.CreatedAt, item.UpdatedAt)
	return err
}

func queryItems(itemType string, limit, offset int) ([]GitHubItem, error) {
	rows, err := db.Query(
		`SELECT id, number, title, state, item_type, author, labels, url, created_at, updated_at
		 FROM github_items WHERE item_type=? ORDER BY updated_at DESC LIMIT ? OFFSET ?`,
		itemType, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []GitHubItem
	for rows.Next() {
		var item GitHubItem
		if err := rows.Scan(&item.ID, &item.Number, &item.Title, &item.State, &item.ItemType, &item.Author, &item.Labels, &item.URL, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func getStats() (Stats, error) {
	var s Stats
	err := db.QueryRow(`SELECT COUNT(*) FROM github_items`).Scan(&s.TotalItems)
	if err != nil {
		return s, err
	}

	db.QueryRow(`SELECT COUNT(*) FROM github_items WHERE item_type='issue' AND state='open'`).Scan(&s.IssuesOpen)
	db.QueryRow(`SELECT COUNT(*) FROM github_items WHERE item_type='issue' AND state='closed'`).Scan(&s.IssuesClosed)
	db.QueryRow(`SELECT COUNT(*) FROM github_items WHERE item_type='pr' AND state='open'`).Scan(&s.PRsOpen)
	db.QueryRow(`SELECT COUNT(*) FROM github_items WHERE item_type='pr' AND (state='closed' OR state='merged')`).Scan(&s.PRsClosed)
	db.QueryRow(`SELECT COUNT(*) FROM github_items WHERE item_type='pr' AND state='merged'`).Scan(&s.PRsMerged)

	return s, nil
}

func logStats() {
	s, err := getStats()
	if err != nil {
		log.Printf("[INFO] Stats query failed: %v\n", err)
		return
	}
	log.Printf("[INFO] GitHub Stats: issues(open=%d, closed=%d) prs(open=%d, closed=%d, merged=%d) total=%d\n",
		s.IssuesOpen, s.IssuesClosed, s.PRsOpen, s.PRsClosed, s.PRsMerged, s.TotalItems)
	appendLog(fmt.Sprintf("Stats: issues(%d/%d) prs(%d/%d/%d)",
		s.IssuesOpen, s.IssuesClosed, s.PRsOpen, s.PRsClosed, s.PRsMerged), s.TotalItems)
}

func saveSnapshot() error {
	s, err := getStats()
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		INSERT INTO snapshots (snapshot_at, issues_open, issues_closed, prs_open, prs_closed, prs_merged, total_items)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, time.Now(), s.IssuesOpen, s.IssuesClosed, s.PRsOpen, s.PRsClosed, s.PRsMerged, s.TotalItems)
	return err
}

func getTrends(limit int) ([]Snapshot, error) {
	rows, err := db.Query(
		`SELECT id, snapshot_at, issues_open, issues_closed, prs_open, prs_closed, prs_merged, total_items
		 FROM snapshots ORDER BY snapshot_at DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []Snapshot
	for rows.Next() {
		var sn Snapshot
		if err := rows.Scan(&sn.ID, &sn.SnapshotAt, &sn.IssuesOpen, &sn.IssuesClosed, &sn.PRsOpen, &sn.PRsClosed, &sn.PRsMerged, &sn.TotalItems); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, sn)
	}
	return snapshots, rows.Err()
}
