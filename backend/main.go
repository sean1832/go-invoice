package main

import (
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"go-invoice/internal/api"
	"go-invoice/internal/auth"
	"go-invoice/internal/crypto"
	"go-invoice/internal/storage"
	"go-invoice/internal/types"
	"go-invoice/internal/ui"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Version is set via ldflags during build
var Version = "dev"

func init() {
	gob.Register(types.UserSessionData{})
}

func main() {
	// Load environment and flags
	loadEnv()
	devmodePtr := flag.Bool("dev", false, "Enable dev mode (uses DEV_FRONTEND_BASE_URL)")
	dbPtr := flag.String("db", "", "Path to the database file.")
	flag.Parse()
	isDevMode := *devmodePtr
	dbPathFromFlag := *dbPtr

	// Load core application configuration
	port, publicURL, frontendURL, storagePath, err := loadAppConfig(isDevMode, dbPathFromFlag)
	if err != nil {
		slog.Error("Failed to load app configuration", "error", err)
		os.Exit(1)
	}

	// Set up session
	sessionConfig, err := setupSession()
	if err != nil {
		slog.Error("Failed to set up session", "error", err)
		os.Exit(1)
	}

	// Set up authentication
	authMethod, err := setupAuth(publicURL, sessionConfig)
	if err != nil {
		slog.Error("Failed to set up authentication", "error", err)
		os.Exit(1)
	}

	// Create the HTTP router (mux)
	mux := http.NewServeMux()

	// Initialize storage
	storageDir, err := storage.NewStorageDir(storagePath)
	if err != nil {
		slog.Error("Failed to initialize storage", "error", err)
		os.Exit(1)
	}

	// Initialize API handler
	apiHandler := api.Handler{
		Context:         context.Background(),
		StorageDir:      *storageDir,
		FrontendBaseURL: frontendURL,
		EmailAuthMethod: authMethod,
		Version:         Version,
	}
	apiHandler.RegisterRoutesV1(mux)

	// Initialize embedded UI handler
	uiHandler, err := ui.NewHandler()
	if err != nil {
		slog.Error("Failed to initialize UI handler", "error", err)
		os.Exit(1)
	}
	mux.Handle("/", uiHandler)

	// Log final config and start the server
	slog.Info("Configuration loaded and services initialized",
		"version", Version,
		"port", port,
		"public_url", publicURL,
		"frontend_url", frontendURL,
		"dev_mode", isDevMode,
		"storage_path", storagePath,
	)

	listenAddr := fmt.Sprintf(":%d", port)
	slog.Info("Server starting...",
		"listen_addr", listenAddr,
		"public_url", publicURL,
	)

	corsHandler := api.WithCORS(mux, []string{frontendURL, "http://localhost:8080", "http://127.0.0.1:8080", "http://0.0.0.0:8080"})
	if err := http.ListenAndServe(listenAddr, corsHandler); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

// loadEnv handles loading the .env file.
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found, relying on OS environment variables.")
	} else {
		slog.Info(".env file loaded.")
	}
}

// loadAppConfig consolidates the loading of port, URLs, and paths.
func loadAppConfig(isDevMode bool, dbPathFromFlag string) (
	port int, publicURL, frontendURL, storagePath string, err error,
) {
	// Port Configuration
	portStr := os.Getenv("PORT")
	port, err = strconv.Atoi(portStr)
	if err != nil || portStr == "" {
		slog.Info("PORT not set or invalid, defaulting to 8080")
		port = 8080
	}

	// Public URL Configuration
	publicURL = os.Getenv("PUBLIC_URL")
	if publicURL == "" {
		publicURL = fmt.Sprintf("http://localhost:%d", port)
		slog.Warn("PUBLIC_URL not set, defaulting to app port", "url", publicURL)
	}
	publicURL = strings.TrimRight(publicURL, "/")

	// Storage Path
	if dbPathFromFlag != "" {
		os.Setenv("STORAGE_PATH", dbPathFromFlag)
		slog.Info("Using custom database path from flag.", "db_path", dbPathFromFlag)
	}
	storagePath = os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		var exePath string
		exePath, err = os.Executable()
		if err != nil {
			err = fmt.Errorf("failed to get executable path: %w", err)
			return
		}
		storagePath = filepath.Join(filepath.Dir(exePath), "db")
		slog.Info("STORAGE_PATH not set, defaulting to executable directory", "storage_path", storagePath)
	}

	// Frontend Base URL
	if isDevMode {
		slog.Info("Dev mode enabled.")
		frontendURL = os.Getenv("DEV_FRONTEND_BASE_URL")
		if frontendURL == "" {
			slog.Warn("DEV_FRONTEND_BASE_URL not set, defaulting to http://localhost:5173")
			frontendURL = "http://localhost:5173"
		}
	} else {
		slog.Info("Production mode enabled.")
		frontendURL = publicURL
	}

	return port, publicURL, frontendURL, storagePath, nil
}

// setupAuth determines and configures the authentication method.
func setupAuth(publicURL string, sessionConfig *auth.SessionConfig) (auth.AuthMethod, error) {
	var authMethod auth.AuthMethod

	// Read auth-specific env vars here
	googleOAuthClientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	googleOAuthClientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	googleAuthAvailable := googleOAuthClientID != "" && googleOAuthClientSecret != ""
	smtpAuthAvailable := smtpPassword != ""

	if !googleAuthAvailable && !smtpAuthAvailable {
		slog.Warn("Neither Google OAuth nor SMTP credentials are set. Email functionality will be disabled.")
		authMethod = auth.AuthMethodNone
	} else if smtpAuthAvailable {
		slog.Warn("Using SMTP_PASSWORD is not recommended for security reasons. Prefer Google OAuth.")
		authMethod = auth.AuthMethodPlain
	} else {
		// Google OAuth is available and preferred
		slog.Info("Google OAuth credentials loaded.")
		authMethod = auth.AuthMethodOAuth2

		callbackURL := fmt.Sprintf("%s/api/v1/mailer/auth/google/callback", publicURL)

		auth.ConfigureGoogleOAuth2(
			googleOAuthClientID,
			googleOAuthClientSecret,
			callbackURL,
			*sessionConfig,
		)
	}

	return authMethod, nil
}

// setupSession configures session parameters from the environment.
func setupSession() (*auth.SessionConfig, error) {
	var key []byte
	keyStr := os.Getenv("SESSION_SECRET")
	if keyStr == "" {
		slog.Info("SESSION_SECRET not set, generating a temporary random key.")
		var err error
		key, err = crypto.GenerateSecureBytes(32)
		if err != nil {
			return nil, fmt.Errorf("error generating secure bytes: %w", err)
		}
	} else {
		key = []byte(keyStr)
	}

	var maxAge int
	ageStr := os.Getenv("SESSION_MAX_AGE")
	maxAge, err := strconv.Atoi(ageStr)
	if ageStr == "" || err != nil {
		maxAge = 86400 * 30 // 30 days default
	}

	var isProd bool
	prodStr := strings.ToLower(os.Getenv("IS_PROD"))
	if prodStr == "true" {
		isProd = true
	}

	// Cookie domain for subdomain support (e.g., ".ztecs.com")
	domain := os.Getenv("COOKIE_DOMAIN")

	return &auth.SessionConfig{
		Key:    key,
		MaxAge: maxAge,
		IsProd: isProd,
		Domain: domain,
	}, nil
}
