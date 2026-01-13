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

	slog.Debug("OAuth2 begin", "provider", provider, "remote_addr", r.RemoteAddr, "host", r.Host)

	// Set provider in query for gothic (needed for session state)
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	// Get state using Gothic's SetState (generates if not present)
	state := gothic.SetState(r)
	slog.Debug("OAuth2 state generated", "state_length", len(state))

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

	slog.Debug("OAuth2 callback received", "provider", provider, "remote_addr", r.RemoteAddr, "host", r.Host)

	// Log cookies for debugging (names only, not values for security)
	var cookieNames []string
	for _, c := range r.Cookies() {
		cookieNames = append(cookieNames, c.Name)
	}
	slog.Debug("OAuth2 callback cookies", "cookie_names", cookieNames)

	// Set provider in query for gothic
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("error complete user auth: %v", err), http.StatusUnauthorized)
		slog.Error("error complete user auth", "error", err, "cookies_present", len(r.Cookies()))
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
