package main

import (
	"context"
	"invoice/internal/api"
	"invoice/internal/storage"
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

	slog.Info("Server listening", "port", port)
	http.ListenAndServe(":"+strconv.Itoa(port), api.WithCORS(mux, []string{"*"}))
}
