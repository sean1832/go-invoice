package main

import (
	"context"
	"flag"
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

	// Initialize logger
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, relying on OS environment variables.")
	} else {
		slog.Info(".env file loaded.")
	}

	// Define and parse flags
	devmodePtr := flag.Bool("dev", false, "Enable dev mode (uses DEV_FRONTEND_BASE_URL)")
	portPtr := flag.Int("port", 8080, "Port for server to host.")
	flag.Parse()

	// Set variables from flags
	port := *portPtr
	isDevMode := *devmodePtr
	var frontendBaseURL string

	// frontend base URL configuration
	if isDevMode {
		// --- Development Mode ---
		slog.Info("Dev mode enabled.")
		frontendBaseURL = os.Getenv("DEV_FRONTEND_BASE_URL")

		if frontendBaseURL == "" {
			slog.Warn("DEV_FRONTEND_BASE_URL not set, defaulting to http://localhost:5173")
			frontendBaseURL = "http://localhost:5173"
		}
	} else {
		// --- Production Mode ---
		slog.Info("Production mode enabled.")
		frontendBaseURL = fmt.Sprintf("http://localhost:%d", port)
	}

	// google oauth client id/secret check
	googleOAuthClientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	googleOAuthClientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	googleOAuthAppPassword := os.Getenv("GOOGLE_OAUTH_APP_PASSWORD")

	if (googleOAuthClientID == "" || googleOAuthClientSecret == "") && googleOAuthAppPassword == "" {
		slog.Warn("GOOGLE_OAUTH_CLIENT_ID, GOOGLE_OAUTH_CLIENT_SECRET, or GOOGLE_OAUTH_APP_PASSWORD must be set in environment variables for email functionality.")
	} else if googleOAuthAppPassword != "" {
		slog.Warn("Using GOOGLE_OAUTH_APP_PASSWORD is not recommended for security reasons.")
	} else {
		slog.Info("Google OAuth credentials loaded.")
	}

	slog.Info("Configuration loaded",
		"port", port,
		"frontend_url", frontendBaseURL,
		"dev_mode", isDevMode,
	)

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

	slog.Info("Server listening", "url", fmt.Sprintf("http://localhost:%d", port))
	http.ListenAndServe(":"+strconv.Itoa(port), api.WithCORS(mux, []string{frontendBaseURL}))
}
