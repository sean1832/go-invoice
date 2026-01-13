package api

import (
	"context"
	"fmt"
	"go-invoice/internal/auth"
	"go-invoice/internal/storage"
	"net/http"
)

const (
	SessionName = "go-invoice-session"
	userKey     = "user"
)

type Handler struct {
	Context         context.Context
	StorageDir      storage.StorageDir
	FrontendBaseURL string
	LocalBaseURL    string // localhost URL for internal PDF generation (ChromeDP)
	EmailAuthMethod auth.AuthMethod
	Version         string
}

func (h *Handler) RegisterRoutesV1(mux *http.ServeMux) {
	const prefix = "/api/v1"
	mux.HandleFunc(prefix+"/", h.Root)
	mux.HandleFunc(prefix+"/version", h.handleVersion)
	mux.HandleFunc(prefix+"/providers", h.handleProvidersCollection)
	mux.HandleFunc(prefix+"/providers/{id}", h.handleProvidersItem)
	mux.HandleFunc(prefix+"/clients", h.handleClientsCollection)
	mux.HandleFunc(prefix+"/clients/{id}", h.handleClientsItem)
	mux.HandleFunc(prefix+"/invoices", h.handleInvoicesCollection)
	mux.HandleFunc(prefix+"/invoices/{id}", h.handleInvoicesItem)
	mux.HandleFunc(prefix+"/invoices/count", h.handleInvoicesCount)
	mux.HandleFunc(prefix+"/invoices/{id}/pdf", h.handleInvoicePDF)
	mux.HandleFunc(fmt.Sprintf("POST %s/invoices/{id}/email", prefix), h.handleSendEmail)
	mux.HandleFunc(fmt.Sprintf("GET %s/email_templates/{id}", prefix), h.handleEmailTemplate)

	// mailer
	mux.HandleFunc(prefix+"/mailer/auth/{provider}", h.handleMailerOAuth2Begin)
	mux.HandleFunc(prefix+"/mailer/auth/{provider}/callback", h.handleMailerOAuth2Callback)
	mux.HandleFunc(fmt.Sprintf("GET %s/mailer/session", prefix), h.handleMailerSession)
	mux.HandleFunc(fmt.Sprintf("POST %s/mailer/logout", prefix), h.handleMailerLogout)
}

func (h *Handler) Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is operational"))
}
