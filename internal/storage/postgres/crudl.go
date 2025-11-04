package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"

	"github.com/P3rCh1/subs-agregator/internal/config"
	"github.com/P3rCh1/subs-agregator/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var ErrNotFound = errors.New("not found")

type SubsAPI interface {
	io.Closer
	Create(ctx context.Context, sub *models.Subscription) error
	Read(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error)
	Summary(ctx context.Context, req *models.SumRequest) (int, error)
}

type subsDB struct {
	db *sqlx.DB
}

func NewSubsAPI(cfg *config.Postgres) (SubsAPI, error) {
	info := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB, cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", info)
	if err != nil {
		return nil, fmt.Errorf("connect postgres fail: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres fail: %w", err)
	}

	return &subsDB{db}, nil
}

func (s *subsDB) Close() error {
	return s.db.Close()
}

func (s *subsDB) Create(ctx context.Context, sub *models.Subscription) error {
	const query = `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	if err := s.db.QueryRowContext(
		ctx,
		query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate,
	).Scan(&sub.ID); err != nil {
		return fmt.Errorf("insert sub fail: %w", err)
	}

	return nil
}

func (s *subsDB) Read(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	const query = `
		SELECT * FROM subscriptions
		WHERE id = $1
	`

	var sub models.Subscription
	if err := s.db.GetContext(ctx, &sub, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("read sub fail: %w", err)
	}

	return &sub, nil
}

func (s *subsDB) Update(ctx context.Context, sub *models.Subscription) error {
	const query = `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6
	`

	res, err := s.db.ExecContext(
		ctx,
		query,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID,
	)

	if err != nil {
		return fmt.Errorf("update sub fail: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *subsDB) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM subscriptions WHERE id = $1
	`

	res, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("delete sub fail: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *subsDB) List(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error) {
	const query = `
		SELECT * FROM subscriptions
		WHERE user_id = $1
	`
	subs := []models.Subscription{}
	if err := s.db.SelectContext(ctx, &subs, query, userID); err != nil {
		return nil, fmt.Errorf("list subs fail: %w", err)
	}

	return subs, nil
}
