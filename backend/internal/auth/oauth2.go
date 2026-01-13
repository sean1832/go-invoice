package auth

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/oauth2"
)

type SessionConfig struct {
	Key    []byte
	MaxAge int
	IsProd bool
	Domain string // Cookie domain (e.g., ".ztecs.com" for subdomains)
}

var GoogleOAuthConfig *oauth2.Config

func ConfigureGoogleOAuth2(clientId, clientSecret, callbackURL string, config SessionConfig) {
	store := sessions.NewCookieStore(config.Key) // TODO: use file store or redis store for production
	store.MaxAge(config.MaxAge)

	// Use SameSite=Lax for OAuth flows - it allows cookies to be sent on top-level navigations
	// SameSite=None requires Secure and can be blocked by some browsers/proxies
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   config.MaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   config.IsProd,
		Domain:   config.Domain, // Set domain for subdomain cookie sharing
	}
	gothic.Store = store

	slog.Info("OAuth2 cookie settings configured",
		"secure", config.IsProd,
		"domain", config.Domain,
		"samesite", "Lax",
		"callback_url", callbackURL,
	)

	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  callbackURL,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"profile",
			"email",
			"https://mail.google.com/", // SMTP access scope (required for OAuth2 SMTP)
		},
	}

	goth.UseProviders(
		google.New(
			clientId, clientSecret, callbackURL,
			// scopes
			"profile",
			"email",
			"https://mail.google.com/", // SMTP access scope (required for OAuth2 SMTP)
		),
	)
}
