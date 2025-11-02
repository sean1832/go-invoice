package main

import (
	"context"
	"invoice/internal/api"
	"invoice/internal/storage"
	"invoice/internal/ui"
	"log/slog"
	"net/http"
	"strconv"
)

func main() {
	port := 8080

	mux := http.NewServeMux()
	storageDir, err := storage.NewStorageDir()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	handler := api.Handler{
		Context:    context.Background(),
		StorageDir: *storageDir,
	}

	handler.RegisterRoutesV1(mux)

	// Serve embedded UI
	uiHandler, err := ui.NewHandler()
	if err != nil {
		slog.Error("Failed to initialize UI handler", "error", err)
		return
	}
	mux.Handle("/", uiHandler)

	slog.Info("Server listening", "port", port)
	http.ListenAndServe(":"+strconv.Itoa(port), api.WithCORS(mux, []string{"*"}))
}
