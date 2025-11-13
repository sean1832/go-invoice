package api

import (
	"go-invoice/internal/auth"
	"net/http"
)

type mailerAuthConfigResponse struct {
	Method auth.AuthMethod `json:"method"`
}

func (h *Handler) handleMailerAuthConfig(w http.ResponseWriter, r *http.Request) {
	method := h.EmailAuthMethod
	resp := mailerAuthConfigResponse{Method: method}
	writeRespOk(w, "auth method obtained", resp)
}
