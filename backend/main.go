package main

import (
	"context"
	"flag"
	"fmt"
	"go-invoice/internal/api"
	"go-invoice/internal/auth"
	"go-invoice/internal/storage"
	"go-invoice/internal/ui"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Version is set via ldflags during build
var Version = "dev"

func main() {

	// Initialize logger
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, relying on OS environment variables.")
	} else {
		slog.Info(".env file loaded.")
	}

	// Define and parse flags
	devmodePtr := flag.Bool("dev", false, "Enable dev mode (uses DEV_FRONTEND_BASE_URL)")
	dbPtr := flag.String("db", "", "Path to the database file.")
	portPtr := flag.Int("port", 8080, "Port for server to host.")
	flag.Parse()

	// Set variables from flags
	port := *portPtr
	isDevMode := *devmodePtr
	dbDir := *dbPtr
	if dbDir != "" {
		os.Setenv("STORAGE_PATH", dbDir)
		slog.Info("Using custom database path from flag.", "db_path", dbDir)
	}
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
	smtp_password := os.Getenv("SMTP_PASSWORD")

	var authMethod auth.AuthMethod
	if (googleOAuthClientID == "" || googleOAuthClientSecret == "") && smtp_password == "" {
		slog.Warn("either GOOGLE_OAUTH_CLIENT_ID, GOOGLE_OAUTH_CLIENT_SECRET, or SMTP_PASSWORD must be set in environment variables for email functionality.")
		authMethod = auth.AuthMethodNone
	} else if smtp_password != "" {
		slog.Warn("Using SMTP_PASSWORD is not recommended for security reasons.")
		authMethod = auth.AuthMethodPlain
	} else {
		slog.Info("Google OAuth credentials loaded.")
		authMethod = auth.AuthMethodOAuth2
	}

	// storage path
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		exePath, err := os.Executable()
		if err != nil {
			slog.Error("Failed to get executable path", "error", err)
			return
		}
		storagePath = filepath.Join(filepath.Dir(exePath), "db")
		slog.Info("STORAGE_PATH not set, defaulting to executable directory", "storage_path", storagePath)
	}

	slog.Info("Configuration loaded",
		"version", Version,
		"port", port,
		"frontend_url", frontendBaseURL,
		"dev_mode", isDevMode,
		"storage_path", storagePath,
	)

	mux := http.NewServeMux()
	storageDir, err := storage.NewStorageDir(storagePath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	handler := api.Handler{
		Context:         context.Background(),
		StorageDir:      *storageDir,
		FrontendBaseURL: frontendBaseURL,
		EmailAuthMethod: authMethod,
		Version:         Version,
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
