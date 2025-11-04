package types

import (
	"fmt"
	"strings"
	"time"
)

const isoDateLayout = "2006-01-02"

// Date wraps time.Time to represent a date as YYYY-MM-DD without time component
type Date struct {
	time.Time
}

// NewDate creates a new Date instance with time set to midnight UTC
func NewDate(t time.Time) Date {
	year, month, day := t.Date()
	return Date{Time: time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// Today returns today's date as a Date instance
func Today() Date {
	return NewDate(time.Now())
}

// AddDate adds years, months, and days to the Date and returns a new Date instance
func (d *Date) AddDate(years int, months int, days int) Date {
	newTime := d.Time.AddDate(years, months, days)
	return NewDate(newTime)
}

// UnmarshalJSON implements the json.Unmarshaler interface for Date
func (d *Date) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	if s == "null" || s == "" {
		d.Time = time.Time{}
		return nil
	}
	t, err := time.Parse(isoDateLayout, s)
	if err != nil {
		// try parse as full RFC3339 timestamp and extract date as fallback
		t, errRFC3339 := time.Parse(time.RFC3339, s)
		if errRFC3339 != nil {
			return fmt.Errorf("failed to parse date as YYYY-MM-DD: %v", errRFC3339)
		}
		*d = NewDate(t.UTC())
		return nil
	}
	*d = Date{Time: t}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Date
func (d *Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", d.Time.Format(isoDateLayout))), nil
}

// String returns the date as a string in YYYY-MM-DD format
func (d *Date) String() string {
	return d.Time.Format(isoDateLayout)
}
