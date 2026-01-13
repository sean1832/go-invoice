package api

import (
	"fmt"
	"go-invoice/internal/auth"
	"go-invoice/internal/types"
	"log/slog"
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"golang.org/x/oauth2"
)

func (h *Handler) handleMailerOAuth2Begin(w http.ResponseWriter, r *http.Request) {
	// Extract provider from path parameter
	provider := r.PathValue("provider")
	if provider == "" {
		writeRespErr(w, "provider not specified", http.StatusBadRequest)
		return
	}

	slog.Info("OAuth2 begin", "provider", provider, "host", r.Host)

	// Forcefully clear ALL existing _gothic_session cookies to prevent stale/expired cookie issues
	// This ensures each login attempt starts completely fresh
	gothicCookieCount := 0
	for _, c := range r.Cookies() {
		if c.Name == "_gothic_session" {
			gothicCookieCount++
		}
	}

	if gothicCookieCount > 0 {
		slog.Info("Clearing all existing _gothic_session cookies", "count", gothicCookieCount)

		// Clear with multiple domain/path combinations to ensure all variations are deleted
		clearVariations := []struct{ domain, path string }{
			{"", "/"},                   // Default
			{"invoice.ztecs.com", "/"},  // Explicit domain
			{".invoice.ztecs.com", "/"}, // Subdomain wildcard
		}

		for _, v := range clearVariations {
			clearCookie := &http.Cookie{
				Name:     "_gothic_session",
				Value:    "",
				Path:     v.path,
				Domain:   v.domain,
				MaxAge:   -1,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, clearCookie)
		}
	}

	// Set provider in query for gothic (needed for session state)
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	// Get state using Gothic's SetState (generates if not present)
	state := gothic.SetState(r)
	slog.Info("OAuth2 state generated and stored in session", "state_prefix", state[:8]+"...")

	// Manually construct OAuth URL with prompt=select_account to force account selection
	if auth.GoogleOAuthConfig == nil {
		writeRespErr(w, "OAuth2 not configured", http.StatusInternalServerError)
		slog.Error("GoogleOAuthConfig is nil")
		return
	}

	authURL := auth.GoogleOAuthConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("prompt", "select_account"),
		oauth2.AccessTypeOffline,
	)

	// Create a session object with the auth URL
	// This needs to be stored so CompleteUserAuth can find it later
	gothProvider, err := goth.GetProvider(provider)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("error getting provider: %v", err), http.StatusInternalServerError)
		slog.Error("error getting provider", "error", err)
		return
	}

	sess, err := gothProvider.BeginAuth(state)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("error beginning auth: %v", err), http.StatusInternalServerError)
		slog.Error("error beginning auth", "error", err)
		return
	}

	// Store the session
	err = gothic.StoreInSession(provider, sess.Marshal(), r, w)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("error storing session: %v", err), http.StatusInternalServerError)
		slog.Error("error storing session", "error", err)
		return
	}

	// Redirect to our custom auth URL with prompt=select_account
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *Handler) handleMailerOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	// Extract provider from path parameter
	provider := r.PathValue("provider")
	if provider == "" {
		writeRespErr(w, "provider not specified", http.StatusBadRequest)
		return
	}

	slog.Info("OAuth2 callback received", "provider", provider, "host", r.Host)

	// Log cookie count for debugging
	gothicCookieCount := 0
	for _, c := range r.Cookies() {
		if c.Name == "_gothic_session" {
			gothicCookieCount++
			slog.Info("_gothic_session cookie", "value_length", len(c.Value))
		}
	}

	if gothicCookieCount > 1 {
		slog.Warn("Multiple _gothic_session cookies still present", "count", gothicCookieCount)
	}

	// Set provider in query for gothic
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		slog.Error("error complete user auth", "error", err, "gothic_cookies", gothicCookieCount)
		writeRespErr(w, fmt.Sprintf("error complete user auth: %v", err), http.StatusUnauthorized)
		return
	}

	slog.Info("User authenticated", "user", user.Email)

	// Get session
	sessions, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		slog.Error("error getting session", "error", err)
		writeRespErr(w, fmt.Sprintf("error getting session: %v", err), http.StatusInternalServerError)
		return
	}

	// Store user data in session
	sessionData := types.UserSessionData{
		Email:        user.Email,
		Name:         user.Name,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		ExpiresAt:    user.ExpiresAt,
		AvatarURL:    user.AvatarURL,
	}
	sessions.Values[userKey] = sessionData

	// Set Location header BEFORE saving session
	redirectURL := fmt.Sprintf("%s/auth-success.html", h.FrontendBaseURL)
	w.Header().Set("Location", redirectURL)

	// Save session - this writes Set-Cookie headers
	if err := sessions.Save(r, w); err != nil {
		slog.Error("error saving session", "error", err)
		// Can't call writeRespErr here as headers may be written
		http.Error(w, fmt.Sprintf("error saving session: %v", err), http.StatusInternalServerError)
		return
	}

	// Write redirect status code LAST
	w.WriteHeader(http.StatusFound)

	slog.Info("OAuth callback completed", "user", user.Email, "redirect", redirectURL)
}
