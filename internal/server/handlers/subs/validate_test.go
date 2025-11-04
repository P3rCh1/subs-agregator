package subs

import (
	"testing"
	"time"

	"github.com/P3rCh1/subs-agregator/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		modifySub func(*models.Subscription)
		wantErr   error
	}{
		{
			name:      "valid data",
			modifySub: func(s *models.Subscription) {},
			wantErr:   nil,
		},
		{
			name:      "negative price",
			modifySub: func(s *models.Subscription) { s.Price = -100 },
			wantErr:   ErrNegativePrice,
		},
		{
			name:      "missing service name",
			modifySub: func(s *models.Subscription) { s.ServiceName = "" },
			wantErr:   ErrServiceNameRequired,
		},
		{
			name:      "missing user id",
			modifySub: func(s *models.Subscription) { s.UserID = uuid.Nil },
			wantErr:   ErrUserIDRequired,
		},
		{
			name:      "missing start date",
			modifySub: func(s *models.Subscription) { s.StartDate = models.MonthDate{} },
			wantErr:   ErrStartDateRequired,
		},
		{
			name: "end date before start date",
			modifySub: func(s *models.Subscription) {
				s.EndDate = models.MonthDate{
					Time:  time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
					Valid: true,
				}
			},
			wantErr: ErrCmpDates,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sub := defaultSub()
			tc.modifySub(&sub)

			err := ValidateSub(&sub)

			if tc.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.wantErr, err)
			}
		})
	}
}
