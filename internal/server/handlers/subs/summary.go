package subs

import (
	"net/http"

	"github.com/P3rCh1/subs-agregator/internal/models"
	"github.com/labstack/echo"
)

func ValidateSumRequest(sr *models.SumRequest) error {
	if sr.StartDate.IsZero() || sr.EndDate.IsZero() {
		return ErrDatesRequired
	}

	if sr.EndDate.Time.Before(sr.StartDate.Time) {
		return ErrCmpDates
	}

	return nil
}

func (s *ServerAPI) Summary(ctx echo.Context) error {
	var r models.SumRequest
	if err := ctx.Bind(&r); err != nil {
		return ErrBadRequest
	}

	if err := ValidateSumRequest(&r); err != nil {
		return err
	}

	sum, err := s.DB.Summary(ctx.Request().Context(), &r)
	if err != nil {
		s.Logger.Error(
			"database",
			"error", err,
		)
		return ErrInternal
	}

	ctx.JSON(http.StatusOK, map[string]int{
		"summary": sum,
	})
	return nil
}
