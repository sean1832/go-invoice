package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func (h *Handler) handleOAuth2Begin(w http.ResponseWriter, r *http.Request) {
	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		// try to get user without re-authenticating
		// TODO: remove this
		fmt.Fprintf(w, "welcome back, %s", user.Email)
		return
	}
	gothic.BeginAuthHandler(w, r)
}

func (h *Handler) handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("error complete user auth: %v", err), http.StatusUnauthorized)
		slog.Error("error complete user auth", "error", err)
		return
	}
	slog.Info("logged in", "user", user.Email, "access token", user.AccessToken)
	writeRespOk(w, "logged in", nil)
}
