package subs

import (
	"net/http"

	"github.com/P3rCh1/subs-aggregator/internal/models"
	"github.com/labstack/echo/v4"
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

// @Summary Calculate total payments
// @Description Calculates the total amount spent on subscriptions within a date range.
// @Description Both start_date and end_date are required; user_id and service_name are optional filters.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body subs.SummaryRequest true "Summary request parameters"
// @Success 200 {object} subs.SummaryResponse
// @Failure 400 {object} subs.ErrorResponse
// @Failure 500 {object} subs.ErrorResponse
// @Router /subs/summary [post]
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
