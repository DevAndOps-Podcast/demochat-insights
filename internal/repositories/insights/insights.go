package insights

import (
	"context"
	"database/sql"
	"demochat-insights/config"
	"errors"
	"fmt"
	"time"
)

type Repository struct {
	db     *sql.DB
	schema string
}

func New(db *sql.DB, cfg *config.Config) *Repository {
	return &Repository{db: db, schema: cfg.DB.Schema}
}

type Insights struct {
	TotalMessages         int64
	MostActiveUserID      int64
	AverageMessageRate    float64
	FirstMessageTimestamp *int64
	LastMessageTimestamp  *int64
}

func (r *Repository) GetInsights(ctx context.Context) (*Insights, error) {
	var i Insights
	err := r.db.QueryRowContext(
		ctx,
		fmt.Sprintf("SELECT total_messages, most_active_user_id, average_message_rate, first_message_timestamp, last_message_timestamp FROM %s.insights WHERE id = 1", r.schema)).
		Scan(&i.TotalMessages, &i.MostActiveUserID, &i.AverageMessageRate, &i.FirstMessageTimestamp, &i.LastMessageTimestamp)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return &i, nil
	}

	return &i, err
}

func (r *Repository) GetUserMessageCount(ctx context.Context, userID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s.user_activity WHERE user_id = $1", r.schema), userID).Scan(&count)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}

	return count, err
}

func (r *Repository) UpdateInsights(ctx context.Context, tx *sql.Tx, insights *Insights) error {
	_, err := tx.ExecContext(
		ctx,
		fmt.Sprintf("UPDATE %s.insights SET total_messages = $1, most_active_user_id = $2, average_message_rate = $3, first_message_timestamp = $4, last_message_timestamp = $5 WHERE id = 1", r.schema),
		insights.TotalMessages, insights.MostActiveUserID, insights.AverageMessageRate, insights.FirstMessageTimestamp, insights.LastMessageTimestamp)
	return err
}

func (r *Repository) UpdateUserActivity(ctx context.Context, tx *sql.Tx, userID int64) error {
	_, err := tx.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s.user_activity (user_id, timestamp) VALUES ($1, $2)", r.schema), userID, time.Now().Unix())
	return err
}

func (r *Repository) GetMostActiveUserID(ctx context.Context) (int64, error) {
	var userID int64
	err := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT user_id FROM %s.user_activity GROUP BY user_id ORDER BY COUNT(*) DESC LIMIT 1", r.schema)).Scan(&userID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return userID, err
}

func (r *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
