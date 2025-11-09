package api

import (
	"fmt"
	"go-invoice/internal/services"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) handleInvoicePDF(w http.ResponseWriter, r *http.Request) {
	// Implementation for generating and returning invoice PDF
	if r.Method != http.MethodGet {
		writeRespErr(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	if id == "" {
		writeRespErr(w, "invoice ID is required", http.StatusBadRequest)
		slog.Error("invoice ID is requred")
		return
	}

	chromeService, err := services.NewChromeService()
	if err != nil {
		writeRespErr(w, "error creating chrome service", http.StatusInternalServerError)
		slog.Error("error creating chrome service", "error", err)
		return
	}
	defer chromeService.Close() // <- finally close

	url := fmt.Sprintf("%s/invoices/%s/print", h.FrontendBaseURL, id)
	slog.Info("generating pdf", "url", url)
	pdf, err := chromeService.GeneratePDF(url, 30*time.Second, services.PaperSizeA3, id)
	if err != nil {
		writeRespErr(w, "error generating pdf", http.StatusInternalServerError)
		slog.Error("error generating pdf", "error", err)
		return
	}

	filename := fmt.Sprintf("%s.pdf", id)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", strconv.Itoa(len(pdf)))

	_, err = w.Write(pdf)
	if err != nil {
		writeRespErr(w, "error writing pdf to response", http.StatusInternalServerError)
		slog.Error("error writing pdf to response", "error", err)
		return
	}
}
