package api

import (
	"go-invoice/internal/storage"
	"net/http"
)

// handleProvidersItem handle api that reads a specific provider information
func (h *Handler) handleProvidersItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getResourceByID(w, r, h.StorageDir.Providers, "provider", func() ResourceData {
			return &storage.ProviderData{}
		})
	case http.MethodPut:
		updateResourceByID(w, r, h.StorageDir.Providers, "provider", func() ResourceData {
			return &storage.ProviderData{}
		})
	case http.MethodDelete:
		deleteResourceByID(w, r, h.StorageDir.Providers, "provider")
	default:
		writeRespErr(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleProvidersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllResources(w, r, h.StorageDir.Providers, "provider", func(dir string) (any, error) {
			return getAllProfiles[*storage.ProviderData](dir)
		})
	case http.MethodPost:
		createResource(w, r, h.StorageDir.Providers, "provider", func() ResourceData {
			return &storage.ProviderData{}
		})
	default:
		writeRespErr(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
