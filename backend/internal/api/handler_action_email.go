package api

import (
	"context"
	"encoding/json"
	"fmt"
	"go-invoice/internal/auth"
	"go-invoice/internal/invoice"
	"go-invoice/internal/services"
	"go-invoice/internal/types"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/markbates/goth/gothic"
	"golang.org/x/oauth2"
)

func (h *Handler) handleSendEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeRespErr(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	logger := slog.With("url", r.RequestURI, "method", r.Method)

	id := r.PathValue("id")
	var emailMessage *types.EmailMessage
	if err := json.NewDecoder(r.Body).Decode(&emailMessage); err != nil {
		writeRespErr(w, fmt.Sprintf("invalid email message data for '%s': %v", id, err), http.StatusBadRequest)
		logger.Error("invalid resource data", "error", err)
		return
	}

	// retreive address and and origin from env
	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	portStr := strings.TrimSpace(os.Getenv("SMTP_PORT"))

	if host == "" || portStr == "" {
		writeRespErr(w, "incomplete SMTP settings", http.StatusInternalServerError)
		logger.Error("incomplete SMTP configuration", "error", "either SMTP_HOST or SMTP_PORT is not configured in environment variables")
		return
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		writeRespErr(w, "malformed SMTP settings", http.StatusInternalServerError)
		logger.Error("failed to parse SMTP_PORT to interger")
		return
	}

	var credential string
	var from string
	switch h.EmailAuthMethod {
	case auth.AuthMethodNone:
		writeRespErr(w, "email sending is not configured", http.StatusNotImplemented)
		return
	case auth.AuthMethodPlain:
		from = strings.TrimSpace(os.Getenv("SMTP_FROM"))
		credential = strings.TrimSpace(os.Getenv("SMTP_PASSWORD"))
		if credential == "" || from == "" {
			writeRespErr(w, "incomplete SMTP configuration, either SMTP_FROM or SMTP_PASSWORD is not set", http.StatusInternalServerError)
			logger.Error("incomplete SMTP configuration", "error", "SMTP_PASSWORD is not set")
			return
		}

	case auth.AuthMethodOAuth2:
		session, err := gothic.Store.Get(r, SessionName)
		if err != nil {
			writeRespErr(w, "failed to get session for oauth2 email sending", http.StatusInternalServerError)
			logger.Error("failed to get session for oauth2 email sending", "error", err)
			return
		}

		val := session.Values[userKey]
		sessionData, ok := val.(types.UserSessionData)
		if !ok {
			writeRespErr(w, "Unauthorized: Not logged in", http.StatusUnauthorized)
			logger.Error("unauthorized: not logged in")
			return
		}

		// refresh token
		// rebuild oauth2 token from session data
		storedToken := &oauth2.Token{
			AccessToken:  sessionData.AccessToken,
			RefreshToken: sessionData.RefreshToken,
			Expiry:       sessionData.ExpiresAt,
		}

		ctx := context.Background()
		tokenSource := auth.GoogleOAuthConfig.TokenSource(ctx, storedToken)
		validToken, err := tokenSource.Token()
		if err != nil {
			writeRespErr(w, "failed to refresh auth token", http.StatusUnauthorized)
			logger.Error("failed to refresh oauth token", "error", err)
			return
		}

		if validToken.AccessToken != storedToken.AccessToken {
			logger.Info("OAuth token refreshed, saving new token to session")
			sessionData.AccessToken = validToken.AccessToken
			sessionData.RefreshToken = validToken.RefreshToken
			sessionData.ExpiresAt = validToken.Expiry
			if err := session.Save(r, w); err != nil {
				writeRespErr(w, "failed to save refreshed token", http.StatusInternalServerError)
				logger.Error("failed to save refreshed token", "error", err)
				return
			}
		}

		from = sessionData.Email
		credential = sessionData.AccessToken

		logger.Info("OAuth2 email sending",
			"from", from,
			"host", host,
			"port", port,
			"token_length", len(credential),
			"token_prefix", credential[:20]+"...")
	}

	smtp := services.NewSMTPService(from, host, port, credential, h.EmailAuthMethod)

	// generate pdf attachment
	chrome, err := services.NewChromeService()
	if err != nil {
		writeRespErr(w, "failed to initialize chrome service for pdf generation", http.StatusInternalServerError)
		logger.Error("failed to initialize chrome service", "error", err)
		return
	}
	defer chrome.Close()
	pdfData, err := chrome.GeneratePDF(fmt.Sprintf("%s/invoices/%s/print", h.FrontendBaseURL, id), 10*time.Second, services.PaperSizeA3, id)
	if err != nil {
		writeRespErr(w, "failed to generate pdf attachment", http.StatusInternalServerError)
		logger.Error("failed to generate pdf attachment", "error", err)
		return
	}
	// send email with attachment
	err = smtp.SendWithAttachment(
		emailMessage.To,
		emailMessage.Subject,
		emailMessage.Body,
		fmt.Sprintf("%s.pdf", id),
		pdfData,
		services.AttachmentTypePDF,
	)
	if err != nil {
		writeRespErr(w, "failed to send email", http.StatusInternalServerError)
		logger.Error("failed to send email", "error", err)
		return
	}

	// update sent status
	inv, err := invoice.LoadInvoice(h.StorageDir.Invoices, id)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("failed to load invoice '%s' to update sent status: %v", id, err), http.StatusInternalServerError)
		logger.Error("failed to load invoice to update sent status", "invoice", id, "error", err)
		return
	}
	inv.Status = invoice.StatusSent // <-- mark as sent
	if err := invoice.SaveInvoice(h.StorageDir.Invoices, inv); err != nil {
		writeRespErr(w, fmt.Sprintf("failed to save invoice '%s' to update sent status: %v", id, err), http.StatusInternalServerError)
		logger.Error("failed to save invoice to update sent status", "invoice", id, "error", err)
		return
	}
	logger.Info("invoice sent status updated", "invoice", id)

	// success
	writeRespOk(w, fmt.Sprintf("email sent for invoice '%s'", id), emailMessage)
	logger.Info("invoice successfully sent", "invoice", id, "from", from, "to", emailMessage.To)
}
