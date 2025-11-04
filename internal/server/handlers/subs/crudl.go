package subs

import (
	"errors"
	"net/http"

	"github.com/P3rCh1/subs-aggregator/internal/models"
	"github.com/P3rCh1/subs-aggregator/internal/storage/postgres"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ValidateSub(sub *models.Subscription) error {
	if sub.Price < 0 {
		return ErrNegativePrice
	}

	if sub.StartDate.IsZero() {
		return ErrStartDateRequired
	}

	if !sub.EndDate.IsZero() && sub.EndDate.Time.Before(sub.StartDate.Time) {
		return ErrCmpDates
	}

	if sub.UserID == uuid.Nil {
		return ErrUserIDRequired
	}

	if sub.ServiceName == "" {
		return ErrServiceNameRequired
	}

	return nil
}

// @Summary Create subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subs [post]
func (s *ServerAPI) Create(ctx echo.Context) error {
	var sub models.Subscription
	if err := ctx.Bind(&sub); err != nil {
		return ErrBadRequest
	}

	if err := ValidateSub(&sub); err != nil {
		return err
	}

	if err := s.DB.Create(ctx.Request().Context(), &sub); err != nil {
		s.Logger.Error(
			"database",
			"error", err,
		)
		return ErrInternal
	}

	ctx.JSON(http.StatusCreated, &sub)
	return nil
}

// @Summary Get subscription
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /subs/{id} [get]
func (s *ServerAPI) Read(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ErrInvalidID
	}

	sub, err := s.DB.Read(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return ErrSubNotFound
		}

		s.Logger.Error(
			"database",
			"error", err,
		)
		return ErrInternal
	}

	ctx.JSON(http.StatusOK, sub)
	return nil
}

// @Summary Update subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body UpdateSubscriptionRequest true "Updated subscription data"
// @Success 200 {object} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /subs/{id} [put]
func (s *ServerAPI) Update(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ErrInvalidID
	}

	var sub models.Subscription
	if err := ctx.Bind(&sub); err != nil {
		return ErrBadRequest
	}

	sub.ID = id

	err = s.DB.Update(ctx.Request().Context(), &sub)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return ErrSubNotFound
		}

		s.Logger.Error(
			"database",
			"error", err,
		)
		return ErrInternal
	}

	ctx.JSON(http.StatusOK, sub)
	return nil
}

// @Summary List subscriptions
// @Tags subscriptions
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {array} SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Router /subs/list/{id} [get]
func (s *ServerAPI) Delete(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ErrInvalidID
	}

	err = s.DB.Delete(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return ErrSubNotFound
		}

		s.Logger.Error(
			"database",
			"error", err,
		)
		return ErrInternal
	}

	ctx.Response().WriteHeader(http.StatusOK)
	return nil
}

// @Summary List subscriptions
// @Tags subscriptions
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {array} models.Subscription
// @Failure 400 {object} map[string]string
// @Router /subs/list/{id} [get]
func (s *ServerAPI) List(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ErrInvalidID
	}

	subs, err := s.DB.List(ctx.Request().Context(), id)
	if err != nil {
		s.Logger.Error(
			"database",
			"error", err,
		)
		return ErrInternal
	}

	ctx.JSON(http.StatusOK, subs)
	return nil
}
