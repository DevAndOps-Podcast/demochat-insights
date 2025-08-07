package repositories

import (
	"database/sql"
	"demochat-insights/config"
	"demochat-insights/internal/repositories/insights"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(func(db *sql.DB, cfg *config.Config) *insights.Repository {
		return insights.New(db, cfg)
	}),
)

