package web

import (
	bots "BotMother/bots/first"
	"BotMother/config"
	"BotMother/logger"
	middlewares "BotMother/middleware"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func StartServer(enabled_bots []string) {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Logger().Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))

	//e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout:      5 * time.Second,
		ErrorMessage: "Request timeout",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			logger.Logger().Error("timeout error",
				// You can add more fields here as needed
				zap.String("error", err.Error()),
				zap.String("method", c.Request().Method),
				zap.String("uri", c.Request().RequestURI),
			)
			c.JSON(http.StatusGatewayTimeout, map[string]string{
				"error": "Request processing timed out",
			})
		},
	}))

	e.Use(middlewares.CheckSecretToken)

	// Setup routes for each bot
	path := fmt.Sprintf(config.Configs.WebhookUrlsFormat, config.Configs.WebhookUrlSecret, "asdklfnasdkfasdkfoasbbot")
	e.POST(path, bots.HandleBotFirst)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", config.Configs.ServerPort)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logger().Fatal("shutting down the server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Logger().Fatal("error during server shutdown", zap.Error(err))
	}
}
