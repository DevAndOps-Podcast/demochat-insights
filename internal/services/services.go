package services

import (
	"demochat-insights/internal/services/insights"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(insights.New),
)

