package api

import (
	"context"
	"go-invoice/internal/storage"
	"net/http"
)

type Handler struct {
	Context    context.Context
	StorageDir storage.StorageDir
}

func (h *Handler) RegisterRoutesV1(mux *http.ServeMux) {
	const prefix = "/api/v1"
	mux.HandleFunc(prefix+"/", h.Root)
	mux.HandleFunc(prefix+"/providers", h.handleProvidersCollection)
	mux.HandleFunc(prefix+"/providers/{id}", h.handleProvidersItem)
	mux.HandleFunc(prefix+"/clients", h.handleClientsCollection)
	mux.HandleFunc(prefix+"/clients/{id}", h.handleClientsItem)
	mux.HandleFunc(prefix+"/invoices", h.handleInvoicesCollection)
	mux.HandleFunc(prefix+"/invoices/{id}", h.handleInvoicesItem)
	mux.HandleFunc(prefix+"/invoices/count", h.handleInvoicesCount)
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is operational"))
}

// // ROUTES:
// // All routes are prefixed with /api/v1

// // PROVIDER PROFILES
// // A "provider" is a resource.
// GET    /providers       // List all provider profiles
// POST   /providers       // Create a new provider profile
// GET    /providers/{id}  // Get a specific provider profile
// PUT    /providers/{id}  // Update a specific provider profile
// DELETE /providers/{id}  // Delete a provider profile

// // CLIENT PROFILES
// // A "client" is a resource.
// GET    /clients         // List all client profiles
// POST   /clients         // Create a new client profile
// GET    /clients/{id}    // Get a specific client profile
// PUT    /clients/{id}    // Update a specific client profile
// DELETE /clients/{id}    // Delete a client profile

// // INVOICES
// // An "invoice" is a resource.
// GET    /invoices        // List all invoices
// POST   /invoices        // Create a new invoice
// GET    /invoices/{id}   // Get a specific invoice
// PUT    /invoices/{id}   // Update a specific invoice
// DELETE /invoices/{id}   // Delete a specific invoice

// // CONFIGURATION
// // This is a "singleton" resource (there's only one).
// GET    /config/smtp     // Get the current SMTP config
// PUT    /config/smtp     // Update the SMTP config

// // ACTIONS (RPC-style)
// // This is for the action of *sending* the email.
// // It's not REST, but it's a clean "action" endpoint.
// POST   /invoices/{id}/send    // Send an invoice.
//                               // The Svelte client will post the generated PDF data
//                               // and profile info here.

// // QUERY
// GET /invoices/count	   // Get total count of invoices

// GET	  /invoices?q={search} // Search invoices by client name, date, or invoice number
// // Supported Query Parameters:
// // -----------------------------------------------------------------
// // Filtering:
// ?client_id={id}      // Filter by a specific client
// ?provider_id={id}    // Filter by a specific provider
// ?status={status}     // e.g., 'draft', 'sent'
// ?due_date_after={date} // e.g., '2025-01-01'
// ?due_date_before={date}// e.g., '2025-01-31'

// // Search:
// ?q={search_term}     // full-text search for client name, number, etc.

// // Pagination:
// ?limit=50            // Number of results to return (default 50)

// // Sorting:
// ?sort=due_date       // Field to sort by (e.g., 'due_date', 'created_at')
// ?order=asc           // 'asc' or 'desc' (default 'desc')
