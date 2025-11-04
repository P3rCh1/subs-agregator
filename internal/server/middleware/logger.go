package middleware

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func Logger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()

			err := next(ctx)

			duration := time.Since(start)

			attributes := []any{
				"method", ctx.Request().Method,
				"path", ctx.Path(),
				"status", ctx.Response().Status,
				"duration", duration.String(),
				"ip", ctx.RealIP(),
				"user_agent", ctx.Request().UserAgent(),
			}

			if err != nil {
				attributes = append(attributes, "error", err.Error())
			}

			logger.Info("request completed", attributes...)

			return err
		}
	}

}
