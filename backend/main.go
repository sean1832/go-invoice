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

	// args --dev
	devmodePtr := flag.Bool("dev", false, "Enable dev mode (seperate frontend & backend port, do not enable this for production!)")
	portPtr := flag.Int("port", 8080, "Port for server to host.")
	flag.Parse()

	if *devmodePtr {
		err := godotenv.Load(".env.dev")
		if err != nil {
			slog.Warn("devmode enabled but .env.dev is not found")
		} else {
			slog.Info("devmode enabled and .env.dev is loaded")
		}
	}

	port := *portPtr
	var frontendBaseURL string
	if *devmodePtr {
		// devmode
		frontendBaseURL = os.Getenv("FRONTEND_BASE_URL")
		if frontendBaseURL == "" {
			slog.Warn("FRONTEND_BASE_URL not set in .env.dev, defaulting to http://localhost:5173")
			frontendBaseURL = "http://localhost:5173"
		}
	} else {
		// production mode
		frontendBaseURL = os.Getenv("FRONTEND_BASE_URL")
		if frontendBaseURL == "" {
			frontendBaseURL = fmt.Sprintf("http://localhost:%d", port)
		}
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

	slog.Info("Server listening", "url", fmt.Sprintf("http://localhost:%d", port))
	http.ListenAndServe(":"+strconv.Itoa(port), api.WithCORS(mux, []string{"*"}))
}
