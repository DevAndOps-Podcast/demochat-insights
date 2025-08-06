package httpapi

import (
	"demochat-insights/internal/httpapi/handlers/insights"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		insights.New,
	),
)
