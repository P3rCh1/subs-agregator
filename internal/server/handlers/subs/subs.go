package subs

import (
	"errors"
	"net/http"

	"github.com/P3rCh1/subs-agregator/internal/models"
	"github.com/P3rCh1/subs-agregator/internal/storage/postgres"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

func Validate(sub *models.Subscription) error {
	if sub.Price < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "negative price")
	}

	if sub.StartDate.IsZero() {
		return echo.NewHTTPError(http.StatusBadRequest, "start_date is required")
	}

	if sub.EndDate != nil && sub.EndDate.Compare(sub.StartDate) <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "end date should be after start date")
	}

	if sub.UserID == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is required")
	}

	if sub.ServiceName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "service_name is required")
	}

	return nil
}

func (s *ServerAPI) Create(ctx echo.Context) error {
	var sub models.Subscription
	if err := ctx.Bind(&sub); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := Validate(&sub); err != nil {
		return err
	}

	if err := s.DB.Create(ctx.Request().Context(), &sub); err != nil {
		s.Logger.Error(
			"database",
			"error", err,
		)
		return echo.ErrInternalServerError
	}

	ctx.JSON(http.StatusCreated, &sub)
	return nil
}

func (s *ServerAPI) Read(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid subscription id")
	}

	sub, err := s.DB.Read(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "subscription not found")
		}

		s.Logger.Error(
			"database",
			"error", err,
		)
		return echo.ErrInternalServerError
	}

	ctx.JSON(http.StatusOK, sub)
	return nil
}

func (s *ServerAPI) Update(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid subscription id")
	}

	var sub models.Subscription
	if err := ctx.Bind(&sub); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	sub.ID = id

	err = s.DB.Update(ctx.Request().Context(), &sub)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "subscription not found")
		}

		s.Logger.Error(
			"database",
			"error", err,
		)
		return echo.ErrInternalServerError
	}

	ctx.JSON(http.StatusOK, sub)
	return nil
}

func (s *ServerAPI) Delete(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid subscription id")
	}

	err = s.DB.Delete(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "subscription not found")
		}

		s.Logger.Error(
			"database",
			"error", err,
		)
		return echo.ErrInternalServerError
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return nil
}

func (s *ServerAPI) List(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	subs, err := s.DB.List(ctx.Request().Context(), id)
	if err != nil {
		s.Logger.Error(
			"database",
			"error", err,
		)
		return echo.ErrInternalServerError
	}

	ctx.JSON(http.StatusOK, subs)
	return nil
}
