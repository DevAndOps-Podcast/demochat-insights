package database

import (
	"context"
	"database/sql"
	"demochat-insights/config"
	"fmt"

	_ "github.com/lib/pq"
)

func New(ctx context.Context, cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	// Set the schema for the current session
	_, err = db.ExecContext(ctx, fmt.Sprintf("SET search_path TO %s", cfg.DB.Schema))
	if err != nil {
		return nil, err
	}

	if err := CreateSchema(db, cfg.DB.Schema); err != nil {
		return nil, err
	}

	if err := createTables(db, cfg.DB.Schema); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateSchema(db *sql.DB, schemaName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName))
	return err
}

func createTables(db *sql.DB, schemaName string) error {
	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.insights (
			id SERIAL PRIMARY KEY,
			total_messages INTEGER NOT NULL DEFAULT 0,
			most_active_user_id INTEGER NOT NULL DEFAULT 0,
			average_message_rate REAL NOT NULL DEFAULT 0,
			first_message_timestamp INTEGER,
			last_message_timestamp INTEGER
		);

		CREATE TABLE IF NOT EXISTS %s.user_activity (
			user_id INTEGER NOT NULL,
			timestamp INTEGER NOT NULL
		);

		INSERT INTO %s.insights (id) VALUES (1) ON CONFLICT (id) DO NOTHING;
	`, schemaName, schemaName, schemaName))

	return err
}
