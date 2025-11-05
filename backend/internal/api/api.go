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