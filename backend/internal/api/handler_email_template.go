package api

import (
	"errors"
	"fmt"
	"go-invoice/internal/storage"
	"net/http"
	"os"
	"path/filepath"
)

func (h *Handler) handleEmailTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeRespErr(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		writeRespErr(w, "email template ID is required", http.StatusBadRequest)
		return
	}

	filepath := filepath.Join(h.StorageDir.EmailTemplates, id+".json")
	emailTemplate := &storage.EmailTemplate{}
	if err := readJSON(filepath, emailTemplate); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeRespErr(w, fmt.Sprintf("email template not found for '%s'", id), http.StatusNotFound)
		} else {
			writeRespErr(w, fmt.Sprintf("failed to read email template '%s'", id), http.StatusInternalServerError)
		}
		return
	}
	writeRespOk(w, fmt.Sprintf("email template '%s'", id), emailTemplate)

}
