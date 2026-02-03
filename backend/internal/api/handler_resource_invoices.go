package api

import (
	"go-invoice/internal/invoice"
	"go-invoice/internal/query"
	"net/http"
	"os"
	"sort"
)

// PaginatedInvoices represents a paginated response of invoices
type PaginatedInvoices struct {
	Items      []invoice.Invoice `json:"items"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalCount int               `json:"total_count"`
	TotalPages int               `json:"total_pages"`
}

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
		invoices, err := getAllInvoices(h.StorageDir.Invoices, "*.json")
		if err != nil && err != os.ErrNotExist {
			writeRespErr(w, "failed to list invoice informations", http.StatusInternalServerError)
			return
		}

		// Apply filters if any
		if queryParams.HasFilters() {
			invoices = query.FilterInvoices(invoices, queryParams)
		}

		// Sort invoices by date descending (newest first)
		sort.Slice(invoices, func(i, j int) bool {
			return invoices[i].Date.After(invoices[j].Date.Time)
		})

		// Calculate pagination
		totalCount := len(invoices)
		totalPages := (totalCount + queryParams.PageSize - 1) / queryParams.PageSize
		if totalPages < 1 {
			totalPages = 1
		}

		// Ensure page is within bounds
		page := queryParams.Page
		if page > totalPages {
			page = totalPages
		}

		// Apply pagination slice
		start := (page - 1) * queryParams.PageSize
		end := start + queryParams.PageSize
		if start > totalCount {
			start = totalCount
		}
		if end > totalCount {
			end = totalCount
		}

		paginatedItems := invoices[start:end]

		// Return paginated response
		result := PaginatedInvoices{
			Items:      paginatedItems,
			Page:       page,
			PageSize:   queryParams.PageSize,
			TotalCount: totalCount,
			TotalPages: totalPages,
		}

		writeRespOk(w, "list of invoices", result)
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

