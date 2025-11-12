package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func (h *Handler) handleOAuth2Begin(w http.ResponseWriter, r *http.Request) {
	// Extract provider from path parameter
	provider := r.PathValue("provider")
	if provider == "" {
		writeRespErr(w, "provider not specified", http.StatusBadRequest)
		return
	}

	// Set provider in query for gothic
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(w, r)
}

func (h *Handler) handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	// Extract provider from path parameter
	provider := r.PathValue("provider")
	if provider == "" {
		writeRespErr(w, "provider not specified", http.StatusBadRequest)
		return
	}

	// Set provider in query for gothic
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("error complete user auth: %v", err), http.StatusUnauthorized)
		slog.Error("error complete user auth", "error", err)
		return
	}
	slog.Info("logged in", "user", user.Email)
	writeRespOk(w, "logged in", nil)
}
