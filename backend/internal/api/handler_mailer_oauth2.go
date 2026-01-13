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

	// Clear any existing gothic session to prevent stale/expired cookie issues
	// This ensures each login attempt starts fresh
	if session, err := gothic.Store.Get(r, "_gothic_session"); err == nil {
		session.Options.MaxAge = -1                        // Mark for deletion
		session.Values = make(map[interface{}]interface{}) // Clear values
		if err := session.Save(r, w); err != nil {
			slog.Warn("Failed to clear old gothic session", "error", err)
		} else {
			slog.Debug("Cleared existing gothic session")
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

	// Log cookies for debugging (names only, not values for security)
	var cookieNames []string
	for _, c := range r.Cookies() {
		cookieNames = append(cookieNames, c.Name)
	}
	slog.Info("OAuth2 callback cookies present", "cookie_names", cookieNames, "count", len(r.Cookies()))

	// Check if the gothic session cookie exists
	gothicSessionExists := false
	for _, c := range r.Cookies() {
		if c.Name == "_gothic_session" {
			gothicSessionExists = true
			break
		}
	}
	if !gothicSessionExists {
		slog.Error("Gothic session cookie missing - this causes state mismatch",
			"hint", "Check: 1) IS_PROD=true for HTTPS, 2) SESSION_SECRET is set, 3) Reverse proxy forwards cookies")
	}

	// Set provider in query for gothic
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("error complete user auth: %v", err), http.StatusUnauthorized)
		slog.Error("error complete user auth", "error", err, "cookies_present", len(r.Cookies()), "gothic_session_exists", gothicSessionExists)
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
		AvatarURL:    user.AvatarURL,
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
