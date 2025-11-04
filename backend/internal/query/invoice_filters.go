package query

import (
	"go-invoice/internal/invoice"
	"strings"
)

// FilterInvoices applies all active filters to the invoice list
func FilterInvoices(invoices []invoice.Invoice, params *InvoiceQueryParams) []invoice.Invoice {
	if !params.HasFilters() {
		return invoices
	}

	filtered := make([]invoice.Invoice, 0, len(invoices))

	for _, inv := range invoices {
		if matchesFilters(inv, params) {
			filtered = append(filtered, inv)
		}
	}

	return filtered
}

// matchesFilters checks if an invoice matches all active filters
func matchesFilters(inv invoice.Invoice, params *InvoiceQueryParams) bool {
	// Client ID filter
	if params.ClientID != "" && !matchesClientID(inv, params.ClientID) {
		return false
	}

	if params.ProviderID != "" && !matchesProviderID(inv, params.ProviderID) {
		return false
	}

	if params.Status != "" && !matchesStatus(inv, params.Status) {
		return false
	}
	if !params.DueDateFrom.IsZero() && inv.Due.Before(params.DueDateFrom.Time) {
		return false
	}
	if !params.DueDateTo.IsZero() && inv.Due.After(params.DueDateTo.Time) {
		return false
	}
	if !params.DateFrom.IsZero() && inv.Date.Before(params.DateFrom.Time) {
		return false
	}
	if !params.DateTo.IsZero() && inv.Date.After(params.DateTo.Time) {
		return false
	}
	return true
}

// matchesClientID checks if invoice's client name matches the ID
func matchesClientID(inv invoice.Invoice, clientID string) bool {
	// client ID is derived from client name (lowercase with underscores)
	return inv.Client.Name == clientID ||
		normalizeID(inv.Client.Name) == clientID
}

// matchesProviderID checks if invoice's provider name matches the ID
func matchesProviderID(inv invoice.Invoice, providerID string) bool {
	// provider ID is derived from provider name (lowercase with underscores)
	return inv.Provider.Name == providerID ||
		normalizeID(inv.Provider.Name) == providerID
}

// normalizeID converts a name to an ID format (lowercase with underscores)
func normalizeID(name string) string {
	result := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	return result
}

func matchesStatus(inv invoice.Invoice, status string) bool {
	return strings.EqualFold(string(inv.Status), status)
}
