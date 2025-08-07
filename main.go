package main

import (
	"context"
	"database/sql"
	"demochat-insights/config"
	"demochat-insights/database"
	"demochat-insights/internal/httpapi"
	"demochat-insights/internal/httpapi/handlers"
	"demochat-insights/internal/repositories"
	"demochat-insights/internal/services"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Lc       fx.Lifecycle
	Config   *config.Config
	DB       *sql.DB
	Echo     *echo.Echo
	Handlers []handlers.Handler `group:"handlers"`
}

func main() {
	app := fx.New(
		fx.Provide(
			context.Background,
			config.New,
			func(ctx context.Context, cfg *config.Config) (*sql.DB, error) {
				return database.New(ctx, cfg)
			},
			echo.New,
		),
		repositories.Module,
		services.Module,
		httpapi.Module,
		fx.Invoke(
			setLifeCycle,
		),
	)

	app.Run()
}

func setLifeCycle(p Params) {
	p.Lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			p.Echo.Use(middleware.Recover())
			p.Echo.Use(middleware.RequestID())

			for _, h := range p.Handlers {
				h.RegisterRoutes(p.Echo)
			}

			go func() {
				p.Echo.Logger.Fatal(p.Echo.Start(p.Config.Address))
			}()

			return nil
		},

		OnStop: func(ctx context.Context) error {
			if err := p.Echo.Shutdown(ctx); err != nil {
				log.Println(err)
			}

			if err := p.DB.Close(); err != nil {
				log.Println(err)
			}

			return nil
		},
	})
}
