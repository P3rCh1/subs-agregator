package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/P3rCh1/subs-agregator/internal/models"
	"github.com/google/uuid"
)

type sumData struct {
	Price     int              `db:"price"`
	StartDate models.MonthDate `db:"start_date"`
	EndDate   models.MonthDate `db:"end_date"`
}

func (s *SubsAPI) Summary(ctx context.Context, req *models.SumRequest) (int, error) {
	conds := []string{
		"start_date <= $2",
		"(end_date IS NULL OR end_date >= $1)",
	}
	args := []any{req.StartDate.Time, req.EndDate.Time}

	if req.ServiceName != "" {
		conds = append(conds, fmt.Sprintf("service_name = $%d", len(args)+1))
		args = append(args, req.ServiceName)
	}

	if req.UserID != uuid.Nil {
		conds = append(conds, fmt.Sprintf("user_id = $%d", len(args)+1))
		args = append(args, req.UserID)
	}

	query := fmt.Sprintf(`
		SELECT price, start_date, end_date
		FROM subscriptions
		WHERE %s
	`, strings.Join(conds, " AND "))

	var subs []sumData
	if err := s.db.SelectContext(ctx, &subs, query, args...); err != nil {
		return 0, fmt.Errorf("summary fetch fail: %w", err)
	}

	total := 0

	for _, sub := range subs {
		firstPay := req.StartDate
		if sub.StartDate.Valid && sub.StartDate.Time.After(req.StartDate.Time) {
			firstPay = sub.StartDate
		}

		end := req.EndDate
		if sub.EndDate.Valid && sub.EndDate.Time.Before(req.EndDate.Time) {
			end = sub.EndDate
		}

		total += sub.Price * firstPay.MonthsBetween(end)
	}

	return total, nil
}
