package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	_ "modernc.org/sqlite"
)

// Open opens the SQLite database stored in the user's application data
// directory and runs the schema migration.
func Open() (*sql.DB, error) {
	dbPath := filepath.Join(xdg.DataHome, "snap-rq-wails-v3", "app.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("creating data directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrating database: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS http_requests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			url TEXT NOT NULL,
			method TEXT NOT NULL,
			body TEXT,
			request_headers TEXT,
			status_code INTEGER,
			response_id TEXT
		);
		CREATE INDEX IF NOT EXISTS idx_http_requests_name ON http_requests(name);
	`)
	return err
}
