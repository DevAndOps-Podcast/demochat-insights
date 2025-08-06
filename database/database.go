package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func New() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./chat.db")
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS insights (
			id INTEGER PRIMARY KEY,
			total_messages INTEGER NOT NULL DEFAULT 0,
			most_active_user_id INTEGER NOT NULL DEFAULT 0,
			average_message_rate REAL NOT NULL DEFAULT 0,
			first_message_timestamp INTEGER,
			last_message_timestamp INTEGER
		);

		CREATE TABLE IF NOT EXISTS user_activity (
			user_id INTEGER NOT NULL,
			timestamp INTEGER NOT NULL
		);

		INSERT OR IGNORE INTO insights (id) VALUES (1);
	`)

	return err
}

