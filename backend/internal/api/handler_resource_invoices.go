package api

import (
	"go-invoice/internal/invoice"
	"go-invoice/internal/query"
	"net/http"
)

func (h *Handler) handleInvoicesItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getResourceByID(w, r, h.StorageDir.Invoices, InvoiceType, func() ResourceData {
			return &invoice.Invoice{}
		})
	case http.MethodPut:
		updateResourceByID(w, r, h.StorageDir.Invoices, InvoiceType, func() ResourceData {
			return &invoice.Invoice{}
		})
	case http.MethodDelete:
		deleteResourceByID(w, r, h.StorageDir.Invoices, InvoiceType)
	default:
		writeRespErr(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleInvoicesCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// TODO: Optimize filtering for large datasets (sqlite)
		// Parse query parameters
		queryParams := query.ParseInvoiceQuery(r.URL.Query())

		// Get all invoices
		invoices, err := getAllInvoices(h.StorageDir.Invoices)
		if err != nil {
			writeRespErr(w, "failed to list invoice informations", http.StatusInternalServerError)
			return
		}

		// Apply filters if any
		if queryParams.HasFilters() {
			invoices = query.FilterInvoices(invoices, queryParams)
		}

		writeRespOk(w, "list of invoices", invoices)
	case http.MethodPost:
		createResource(w, r, h.StorageDir.Invoices, InvoiceType, func() ResourceData {
			return &invoice.Invoice{}
		})
	default:
		writeRespErr(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleInvoicesCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeRespErr(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	count, err := countResources(h.StorageDir.Invoices)
	if err != nil {
		writeRespErr(w, "failed to count invoices", http.StatusInternalServerError)
		return
	}
	writeRespOk(w, "invoice count retrieved", map[string]int{"count": count})
}
