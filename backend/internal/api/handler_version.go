package api

import (
	"encoding/json"
	"net/http"
)

// VersionResponse represents the version information
type VersionResponse struct {
	Version string `json:"version"`
}

// handleVersion returns the application version
func (h *Handler) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeRespErr(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := VersionResponse{
		Version: h.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		writeRespErr(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
