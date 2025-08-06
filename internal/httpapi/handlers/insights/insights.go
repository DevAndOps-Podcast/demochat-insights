package insights

import (
	"demochat-insights/config"
	"demochat-insights/internal/httpapi/handlers"
	"demochat-insights/internal/services/insights"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

type Handler struct {
	insightsService *insights.Service
	apiKey          string
}

type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

func New(insightsService *insights.Service, cfg *config.Config) Result {
	return Result{
		Handler: &Handler{
			insightsService: insightsService,
			apiKey:          cfg.ApiKey,
		},
	}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	apiKeyMiddleware := middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:service-secret",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == h.apiKey, nil
		},
	})

	group := e.Group("/messages", apiKeyMiddleware)

	group.POST("", h.saveMessage)
	group.GET("", h.getInsights)
}

func (h *Handler) saveMessage(c echo.Context) error {
	log.Println("saving message")
	var req insights.PublishMessageRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.insightsService.PublishMessage(c.Request().Context(), req); err != nil {
		log.Println("failed to save message", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	log.Println("insights updated")
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) getInsights(c echo.Context) error {
	insights, err := h.insightsService.GetInsights(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, insights)
}
