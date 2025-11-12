package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type SessionConfig struct {
	Key    []byte
	MaxAge int
	IsProd bool
}

func ConfigureGoogleOAuth2(clientId, clientSecret, callbackURL string, config SessionConfig) {
	store := sessions.NewCookieStore(config.Key)
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

	goth.UseProviders(
		google.New(
			clientId, clientSecret, callbackURL,
			// scopes
			"profile", "email", // gets user's profile info & email
			"https://www.googleapis.com/auth/gmail.send", // permission to send email
		),
	)
}
