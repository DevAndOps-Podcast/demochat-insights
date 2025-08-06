package repositories

import (
	"demochat-insights/internal/repositories/insights"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(insights.New),
)

