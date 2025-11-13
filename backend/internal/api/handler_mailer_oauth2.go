package api

import (
	"fmt"
	"go-invoice/internal/types"
	"log/slog"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func (h *Handler) handleMailerOAuth2Begin(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) handleMailerOAuth2Callback(w http.ResponseWriter, r *http.Request) {
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

	sessions, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("error getting session: %v", err), http.StatusInternalServerError)
		slog.Error("error getting session", "error", err)
		return
	}
	sessionData := types.UserSessionData{
		Email:        user.Email,
		Name:         user.Name,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		ExpiresAt:    user.ExpiresAt,
	}
	sessions.Values[userKey] = sessionData

	// Save session
	if err := sessions.Save(r, w); err != nil {
		writeRespErr(w, fmt.Sprintf("error saving session: %v", err), http.StatusInternalServerError)
		slog.Error("error saving session", "error", err)
		return
	}

	slog.Info("logged in", "user", user.Email)

	// Redirect to auth success page (will close popup)
	w.Header().Set("Location", fmt.Sprintf("%s/auth-success.html", h.FrontendBaseURL))
	w.WriteHeader(http.StatusFound)
}
