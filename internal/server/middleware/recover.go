package middleware

import (
	"net/http"
	"runtime/debug"

	"log/slog"

	"github.com/labstack/echo/v4"
)

func Recover(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					logger.Error(
						"panic recovered",
						"error", r,
						"path", c.Path(),
						"method", c.Request().Method,
						"stack", string(debug.Stack()),
					)

					c.JSON(http.StatusInternalServerError, map[string]any{
						"message": "internal server error",
					})
				}
			}()

			return next(c)
		}
	}
}
