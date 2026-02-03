package query

import (
	"go-invoice/internal/invoice"
	"testing"
)

func TestPaginateInvoices(t *testing.T) {
	// Create test invoices
	createInvoices := func(count int) []invoice.Invoice {
		invoices := make([]invoice.Invoice, count)
		for i := 0; i < count; i++ {
			invoices[i] = invoice.Invoice{ID: string(rune('A' + i))}
		}
		return invoices
	}

	tests := []struct {
		name           string
		totalItems     int
		page           int
		pageSize       int
		expectedStart  int
		expectedEnd    int
		expectedCount  int
		expectedPages  int
	}{
		{
			name:          "first page of many",
			totalItems:    50,
			page:          1,
			pageSize:      10,
			expectedStart: 0,
			expectedEnd:   10,
			expectedCount: 10,
			expectedPages: 5,
		},
		{
			name:          "middle page",
			totalItems:    50,
			page:          3,
			pageSize:      10,
			expectedStart: 20,
			expectedEnd:   30,
			expectedCount: 10,
			expectedPages: 5,
		},
		{
			name:          "last page partial",
			totalItems:    25,
			page:          3,
			pageSize:      10,
			expectedStart: 20,
			expectedEnd:   25,
			expectedCount: 5,
			expectedPages: 3,
		},
		{
			name:          "page beyond total - clamp to last",
			totalItems:    25,
			page:          10,
			pageSize:      10,
			expectedStart: 20,
			expectedEnd:   25,
			expectedCount: 5,
			expectedPages: 3,
		},
		{
			name:          "empty list",
			totalItems:    0,
			page:          1,
			pageSize:      10,
			expectedStart: 0,
			expectedEnd:   0,
			expectedCount: 0,
			expectedPages: 1,
		},
		{
			name:          "single item",
			totalItems:    1,
			page:          1,
			pageSize:      10,
			expectedStart: 0,
			expectedEnd:   1,
			expectedCount: 1,
			expectedPages: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoices := createInvoices(tt.totalItems)
			totalCount := len(invoices)

			// Calculate pagination (same logic as handler)
			totalPages := (totalCount + tt.pageSize - 1) / tt.pageSize
			if totalPages < 1 {
				totalPages = 1
			}

			page := tt.page
			if page > totalPages {
				page = totalPages
			}

			start := (page - 1) * tt.pageSize
			end := start + tt.pageSize
			if start > totalCount {
				start = totalCount
			}
			if end > totalCount {
				end = totalCount
			}

			paginatedItems := invoices[start:end]

			if len(paginatedItems) != tt.expectedCount {
				t.Errorf("got %d items, want %d", len(paginatedItems), tt.expectedCount)
			}
			if totalPages != tt.expectedPages {
				t.Errorf("got %d pages, want %d", totalPages, tt.expectedPages)
			}
			if start != tt.expectedStart {
				t.Errorf("start = %d, want %d", start, tt.expectedStart)
			}
			if end != tt.expectedEnd {
				t.Errorf("end = %d, want %d", end, tt.expectedEnd)
			}
		})
	}
}
