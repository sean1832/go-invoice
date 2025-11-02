package query

import (
	"net/url"
	"time"
)

// InvoiceQueryParams holds all possible query parameters for invoice filtering
type InvoiceQueryParams struct {
	ClientID    string
	ProviderID  string
	Status      string
	DueDateFrom time.Time
	DueDateTo   time.Time
	DateFrom    time.Time
	DateTo      time.Time
}

// ParseInvoiceQuery extracts and validates query parameters from URL
func ParseInvoiceQuery(values url.Values) *InvoiceQueryParams {
	return &InvoiceQueryParams{
		ClientID:    values.Get("client_id"),
		ProviderID:  values.Get("provider_id"),
		Status:      values.Get("status"),
		DueDateFrom: parseTimeParam(values.Get("due_from")),
		DueDateTo:   parseTimeParam(values.Get("due_to")),
		DateFrom:    parseTimeParam(values.Get("from")),
		DateTo:      parseTimeParam(values.Get("to")),
	}
}

// HasFilters returns true if any filter is active
func (q *InvoiceQueryParams) HasFilters() bool {
	return (q.ClientID != "" ||
		q.ProviderID != "" ||
		q.Status != "" ||
		!q.DueDateFrom.IsZero() ||
		!q.DueDateTo.IsZero() ||
		!q.DateFrom.IsZero() ||
		!q.DateTo.IsZero())
}

func parseTimeParam(value string) time.Time {
	const format = "2006-01-02"
	t, err := time.Parse(format, value)
	if err != nil {
		return time.Time{}
	}
	return t
}
