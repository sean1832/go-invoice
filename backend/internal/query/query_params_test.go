package query

import (
	"net/url"
	"testing"
)

func TestParseInvoiceQuery_Pagination(t *testing.T) {
	tests := []struct {
		name             string
		queryString      string
		expectedPage     int
		expectedPageSize int
	}{
		{
			name:             "default values when no params",
			queryString:      "",
			expectedPage:     DefaultPage,
			expectedPageSize: DefaultPageSize,
		},
		{
			name:             "custom page and page_size",
			queryString:      "page=3&page_size=10",
			expectedPage:     3,
			expectedPageSize: 10,
		},
		{
			name:             "page only - uses default page_size",
			queryString:      "page=5",
			expectedPage:     5,
			expectedPageSize: DefaultPageSize,
		},
		{
			name:             "page_size only - uses default page",
			queryString:      "page_size=50",
			expectedPage:     DefaultPage,
			expectedPageSize: 50,
		},
		{
			name:             "invalid page defaults to 1",
			queryString:      "page=invalid",
			expectedPage:     DefaultPage,
			expectedPageSize: DefaultPageSize,
		},
		{
			name:             "negative page defaults to 1",
			queryString:      "page=-5",
			expectedPage:     DefaultPage,
			expectedPageSize: DefaultPageSize,
		},
		{
			name:             "zero page defaults to 1",
			queryString:      "page=0",
			expectedPage:     DefaultPage,
			expectedPageSize: DefaultPageSize,
		},
		{
			name:             "page_size exceeds max - capped at max",
			queryString:      "page_size=500",
			expectedPage:     DefaultPage,
			expectedPageSize: MaxPageSize,
		},
		{
			name:             "negative page_size defaults to default",
			queryString:      "page_size=-10",
			expectedPage:     DefaultPage,
			expectedPageSize: DefaultPageSize,
		},
		{
			name:             "zero page_size defaults to default",
			queryString:      "page_size=0",
			expectedPage:     DefaultPage,
			expectedPageSize: DefaultPageSize,
		},
		{
			name:             "mixed with filters",
			queryString:      "page=2&page_size=25&status=draft&client_id=test",
			expectedPage:     2,
			expectedPageSize: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, _ := url.ParseQuery(tt.queryString)
			result := ParseInvoiceQuery(values)

			if result.Page != tt.expectedPage {
				t.Errorf("Page = %d, want %d", result.Page, tt.expectedPage)
			}
			if result.PageSize != tt.expectedPageSize {
				t.Errorf("PageSize = %d, want %d", result.PageSize, tt.expectedPageSize)
			}
		})
	}
}

func TestParseInvoiceQuery_Filters(t *testing.T) {
	values, _ := url.ParseQuery("client_id=john_doe&provider_id=company_abc&status=draft")
	result := ParseInvoiceQuery(values)

	if result.ClientID != "john_doe" {
		t.Errorf("ClientID = %s, want john_doe", result.ClientID)
	}
	if result.ProviderID != "company_abc" {
		t.Errorf("ProviderID = %s, want company_abc", result.ProviderID)
	}
	if result.Status != "draft" {
		t.Errorf("Status = %s, want draft", result.Status)
	}
}

func TestHasFilters(t *testing.T) {
	tests := []struct {
		name     string
		params   *InvoiceQueryParams
		expected bool
	}{
		{
			name:     "no filters",
			params:   &InvoiceQueryParams{Page: 1, PageSize: 20},
			expected: false,
		},
		{
			name:     "client_id filter",
			params:   &InvoiceQueryParams{ClientID: "test"},
			expected: true,
		},
		{
			name:     "provider_id filter",
			params:   &InvoiceQueryParams{ProviderID: "test"},
			expected: true,
		},
		{
			name:     "status filter",
			params:   &InvoiceQueryParams{Status: "draft"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.params.HasFilters()
			if result != tt.expected {
				t.Errorf("HasFilters() = %v, want %v", result, tt.expected)
			}
		})
	}
}
