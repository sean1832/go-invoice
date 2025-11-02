package api

import (
	"invoice/internal/invoice"
	"net/http"
)

func (h *Handler) handleInvoicesItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getResourceByID(w, r, h.StorageDir.Invoices, "invoice", func() ResourceData {
			return &invoice.Invoice{}
		})
	case http.MethodPut:
		updateResourceByID(w, r, h.StorageDir.Invoices, "invoice", func() ResourceData {
			return &invoice.Invoice{}
		})
	case http.MethodDelete:
		deleteResourceByID(w, r, h.StorageDir.Invoices, "invoice")
	default:
		writeRespErr(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleInvoicesCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		

		getAllResources(w, h.StorageDir.Invoices, "invoice", func(root string) (any, error) {
			return getAllInvoices(root)
		})
	case http.MethodPost:
		createResource(w, r, h.StorageDir.Invoices, "invoice", func() ResourceData {
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