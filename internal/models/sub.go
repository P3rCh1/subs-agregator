package models

import (
	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `json:"id"           db:"id"`
	ServiceName string    `json:"service_name" db:"service_name"`
	Price       int       `json:"price"        db:"price"`
	UserID      uuid.UUID `json:"user_id"      db:"user_id"`
	StartDate   MonthDate `json:"start_date"   db:"start_date"`
	EndDate     MonthDate `json:"end_date"     db:"end_date"`
}

type SumRequest struct {
	ServiceName string    `json:"service_name"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   MonthDate `json:"start_date"`
	EndDate     MonthDate `json:"end_date"`
}
