package api

import (
	"invoice/internal/storage"
	"net/http"
)

func (h *Handler) handleClientsItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getResourceByID(w, r, h.StorageDir.Clients, "client", func() ResourceData {
			return &storage.ClientData{}
		})
	case http.MethodPut:
		updateResourceByID(w, r, h.StorageDir.Clients, "client", func() ResourceData {
			return &storage.ClientData{}
		})
	case http.MethodDelete:
		deleteResourceByID(w, r, h.StorageDir.Clients, "client")
	default:
		writeRespErr(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleClientsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllResources(w, h.StorageDir.Clients, "client", func(dir string) (any, error) {
			return getAllProfiles[*storage.ClientData](dir)
		})
	case http.MethodPost:
		createResource(w, r, h.StorageDir.Clients, "client", func() ResourceData {
			return &storage.ClientData{}
		})
	default:
		writeRespErr(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
