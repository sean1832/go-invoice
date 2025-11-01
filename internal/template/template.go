package template

import (
	"invoice/internal/invoice"
)

// Data represents the data structure used to render HTML invoice templates.
// This is specifically designed for template rendering and may differ from
// the core Invoice domain model to suit presentation needs.
type Data struct {
	// Re-export all invoice fields for template access
	// This allows the HTML template to use the same field names
	*invoice.Invoice
}

// FromInvoice converts a core Invoice to template Data
func FromInvoice(inv *invoice.Invoice) *Data {
	return &Data{
		Invoice: inv,
	}
}
