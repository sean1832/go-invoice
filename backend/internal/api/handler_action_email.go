package api

import (
	"encoding/json"
	"fmt"
	"go-invoice/internal/auth"
	"go-invoice/internal/services"
	"go-invoice/internal/types"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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

	switch h.EmailAuthMethod {
	case auth.AuthMethodNone:
		writeRespErr(w, "email sending is not configured", http.StatusNotImplemented)
		return
	case auth.AuthMethodPlain:
		password := strings.TrimSpace(os.Getenv("SMTP_PASSWORD"))
		from := strings.TrimSpace(os.Getenv("SMTP_FROM"))
		host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
		portStr := strings.TrimSpace(os.Getenv("SMTP_PORT"))
		if from == "" || host == "" || portStr == "" || password == "" {
			writeRespErr(w, "incomplete SMTP settings", http.StatusInternalServerError)
			logger.Error("incomplete SMTP configuration, either SMTP_FROM, SMTP_HOST, SMTP_PORT, SMTP_PASSWORD is not configured in environment variables")
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			writeRespErr(w, "malformed SMTP settings", http.StatusInternalServerError)
			logger.Error("failed to parse SMTP_PORT to interger")
		}

		smtp := services.NewSMTPService(from, host, port, password)
		// generate pdf attachment
		chrome, err := services.NewChromeService()
		if err != nil {
			writeRespErr(w, "failed to initialize chrome service for pdf generation", http.StatusInternalServerError)
			logger.Error("failed to initialize chrome service", "error", err)
			return
		}
		defer chrome.Close()
		pdfData, err := chrome.GeneratePDF(fmt.Sprintf("%s/invoices/%s/print", h.FrontendBaseURL, id), 10*time.Second)
		if err != nil {
			writeRespErr(w, "failed to generate pdf attachment", http.StatusInternalServerError)
			logger.Error("failed to generate pdf attachment", "error", err)
			return
		}
		// send email with attachment
		smtp.SendWithAttachment(
			emailMessage.To,
			emailMessage.Subject,
			emailMessage.Body,
			fmt.Sprintf("%s.pdf", id),
			pdfData,
			services.AttachmentTypePDF,
		)

	case auth.AuthMethodOAuth2:
		// TODO: Implement email sending via OAuth2
		writeRespErr(w, "email sending via OAuth2 is not implemented", http.StatusNotImplemented)
		return
	}

	writeRespOk(w, fmt.Sprintf("email sent for invoice '%s'", id), emailMessage)
}
