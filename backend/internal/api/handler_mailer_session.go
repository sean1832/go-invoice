package api

import (
	"go-invoice/internal/auth"
	"go-invoice/internal/types"
	"log/slog"
	"net/http"

	"github.com/markbates/goth/gothic"
)

// SessionResponse represents the authentication session status
type SessionResponse struct {
	Authenticated bool   `json:"authenticated"`
	Email         string `json:"email,omitempty"`
	Method        string `json:"method"`
}

// handleMailerSession checks if the user is authenticated and returns session info
// GET /api/v1/mailer/session
func (h *Handler) handleMailerSession(w http.ResponseWriter, r *http.Request) {
	// Get auth method
	authMethod := h.EmailAuthMethod
	methodStr := authMethodToString(authMethod)

	// If not using OAuth2, no session to check
	if authMethod != auth.AuthMethodOAuth2 {
		writeRespOk(w, "session status", SessionResponse{
			Authenticated: false,
			Method:        methodStr,
		})
		return
	}

	// Get session
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		slog.Warn("error getting session", "error", err)
		writeRespOk(w, "session status", SessionResponse{
			Authenticated: false,
			Method:        methodStr,
		})
		return
	}

	// Check if user is in session
	userData, ok := session.Values[userKey].(types.UserSessionData)
	if !ok || userData.Email == "" {
		writeRespOk(w, "session status", SessionResponse{
			Authenticated: false,
			Method:        methodStr,
		})
		return
	}

	// User is authenticated
	writeRespOk(w, "session status", SessionResponse{
		Authenticated: true,
		Email:         userData.Email,
		Method:        methodStr,
	})
}

// handleMailerLogout clears the user session
// POST /api/v1/mailer/logout
func (h *Handler) handleMailerLogout(w http.ResponseWriter, r *http.Request) {
	// Get session
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		slog.Warn("error getting session for logout", "error", err)
		writeRespOk(w, "logged out", map[string]bool{"success": true})
		return
	}

	// Get user email for logging
	userEmail := ""
	if userData, ok := session.Values[userKey].(types.UserSessionData); ok {
		userEmail = userData.Email
	}

	// Clear session
	session.Values[userKey] = types.UserSessionData{}
	session.Options.MaxAge = -1 // Delete cookie

	// Save session (to delete it)
	if err := session.Save(r, w); err != nil {
		writeRespErr(w, "error clearing session", http.StatusInternalServerError)
		slog.Error("error clearing session", "error", err)
		return
	}

	slog.Info("user logged out", "email", userEmail)
	writeRespOk(w, "logged out", map[string]bool{"success": true})
}

// authMethodToString converts auth.AuthMethod to string
func authMethodToString(method auth.AuthMethod) string {
	switch method {
	case auth.AuthMethodOAuth2:
		return "oauth2"
	case auth.AuthMethodPlain:
		return "plain"
	case auth.AuthMethodNone:
		return "none"
	default:
		return "none"
	}
}
