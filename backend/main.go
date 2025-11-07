package main

import (
	"context"
	"fmt"
	"go-invoice/internal/api"
	"go-invoice/internal/storage"
	"go-invoice/internal/ui"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.dev")
	if err == nil {
		slog.Info(".env.dev file found!")
	}
	port := 8080
	frontendBaseURL := os.Getenv("FRONTEND_BASE_URL")
	if frontendBaseURL == "" {
		frontendBaseURL = fmt.Sprintf("http://localhost:%d", port)
	}

	mux := http.NewServeMux()
	storageDir, err := storage.NewStorageDir()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	handler := api.Handler{
		Context:         context.Background(),
		StorageDir:      *storageDir,
		FrontendBaseURL: frontendBaseURL,
	}

	handler.RegisterRoutesV1(mux)

	// Serve embedded UI
	uiHandler, err := ui.NewHandler()
	if err != nil {
		slog.Error("Failed to initialize UI handler", "error", err)
		return
	}
	mux.Handle("/", uiHandler)

	slog.Info("Server listening", "url", "http://localhost:"+strconv.Itoa(port))
	http.ListenAndServe(":"+strconv.Itoa(port), api.WithCORS(mux, []string{"*"}))
}
