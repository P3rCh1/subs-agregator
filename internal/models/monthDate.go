package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type MonthDate struct {
	Time  time.Time
	Valid bool
}

const layout = "01-2006"

func (m *MonthDate) UnmarshalJSON(b []byte) error {
	s := string(b)

	if s == "null" || s == `""` {
		m.Time = time.Time{}
		m.Valid = false
		return nil
	}

	s = s[1 : len(s)-1]

	t, err := time.Parse(layout, s)
	if err != nil {
		return fmt.Errorf("invalid date format (expected MM-YYYY): %w", err)
	}

	m.Time = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	m.Valid = true
	return nil
}

func (m MonthDate) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%02d-%04d\"", m.Time.Month(), m.Time.Year())), nil
}

func (m *MonthDate) Scan(v any) error {
	if v == nil {
		m.Time = time.Time{}
		m.Valid = false
		return nil
	}

	t, ok := v.(time.Time)
	if !ok {
		return fmt.Errorf("MonthDate.Scan: invalid DB value %T", v)
	}

	m.Time = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	m.Valid = true
	return nil
}

func (m MonthDate) Value() (driver.Value, error) {
	if !m.Valid {
		return nil, nil
	}
	return m.Time, nil
}

func (m MonthDate) IsZero() bool {
	return !m.Valid || m.Time.IsZero()
}

func (m MonthDate) MonthsBetween(other MonthDate) int {
	if !m.Valid || !other.Valid {
		return 0
	}

	y1, m1, _ := m.Time.Date()
	y2, m2, _ := other.Time.Date()

	return (y2-y1)*12 + int(m2-m1) + 1
}
