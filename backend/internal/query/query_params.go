package query

import (
	"go-invoice/internal/types"
	"net/url"
	"strconv"
	"time"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// InvoiceQueryParams holds all possible query parameters for invoice filtering
type InvoiceQueryParams struct {
	ClientID    string
	ProviderID  string
	Status      string
	DueDateFrom types.Date
	DueDateTo   types.Date
	DateFrom    types.Date
	DateTo      types.Date
	// Pagination
	Page     int
	PageSize int
}

// ParseInvoiceQuery extracts and validates query parameters from URL
func ParseInvoiceQuery(values url.Values) *InvoiceQueryParams {
	page := parseIntParam(values.Get("page"), DefaultPage)
	pageSize := parseIntParam(values.Get("page_size"), DefaultPageSize)

	// Validate page
	if page < 1 {
		page = DefaultPage
	}

	// Validate page_size
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return &InvoiceQueryParams{
		ClientID:    values.Get("client_id"),
		ProviderID:  values.Get("provider_id"),
		Status:      values.Get("status"),
		DueDateFrom: parseTimeParam(values.Get("due_from")),
		DueDateTo:   parseTimeParam(values.Get("due_to")),
		DateFrom:    parseTimeParam(values.Get("from")),
		DateTo:      parseTimeParam(values.Get("to")),
		Page:        page,
		PageSize:    pageSize,
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

func parseTimeParam(value string) types.Date {
	const format = "2006-01-02"
	t, err := time.Parse(format, value)
	if err != nil {
		return types.Date{}
	}
	return types.NewDate(t)
}

func parseIntParam(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

