package subs

import (
	"log/slog"
	"net/http"

	"github.com/P3rCh1/subs-agregator/internal/config"
	"github.com/P3rCh1/subs-agregator/internal/storage/postgres"
	"github.com/labstack/echo"
)

var (
	ErrInternal            = echo.NewHTTPError(http.StatusBadRequest, "internal server error")
	ErrBadRequest          = echo.NewHTTPError(http.StatusBadRequest, "bad request")
	ErrNegativePrice       = echo.NewHTTPError(http.StatusBadRequest, "negative price")
	ErrStartDateRequired   = echo.NewHTTPError(http.StatusBadRequest, "start_date is required")
	ErrDatesRequired       = echo.NewHTTPError(http.StatusBadRequest, "dates are required")
	ErrCmpDates            = echo.NewHTTPError(http.StatusBadRequest, "end date should be after start date")
	ErrServiceNameRequired = echo.NewHTTPError(http.StatusBadRequest, "service_name is required")
	ErrUserIDRequired      = echo.NewHTTPError(http.StatusBadRequest, "user_id is required")
	ErrInvalidID           = echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	ErrSubNotFound         = echo.NewHTTPError(http.StatusBadRequest, "subscription not found")
)

type ServerAPI struct {
	Logger *slog.Logger
	Config *config.Config
	DB     *postgres.SubsAPI
}

func NewServerAPI(logger *slog.Logger, config *config.Config, db *postgres.SubsAPI) *ServerAPI {
	return &ServerAPI{
		Logger: logger,
		Config: config,
		DB:     db,
	}
}
