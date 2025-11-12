package auth

import (
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
}

var GoogleOAuthConfig *oauth2.Config

func ConfigureGoogleOAuth2(clientId, clientSecret, callbackURL string, config SessionConfig) {
	store := sessions.NewCookieStore(config.Key) // TODO: use file store or redis store for production
	store.MaxAge(config.MaxAge)

	var sameSite http.SameSite
	if config.IsProd {
		sameSite = http.SameSiteNoneMode
	} else {
		sameSite = http.SameSiteLaxMode
	}

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   config.MaxAge,
		HttpOnly: true,
		SameSite: sameSite,
		Secure:   config.IsProd,
	}
	gothic.Store = store

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
