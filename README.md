# go-invoice

Simple no bullshit invoice management application built with SvelteKit and Go.

> [!WARNING]
> This project is still work in progress. Some features may be incomplete or missing.

## Features

- Create, view, and manage invoices
- Client and provider management
- REST API for integration

## Quick Start

```bash
# Build the application
npm run build

# Run the application
./backend/bin/invoice-app      # Linux/Mac
.\backend\bin\invoice-app.exe  # Windows
```

The server will start at http://localhost:8080

## Development

```bash
# Install frontend dependencies
npm run install:frontend

# Run in development mode
npm run dev:frontend  # Terminal 1 - http://localhost:5173
npm run dev:backend   # Terminal 2 - http://localhost:8080
```

Or use VS Code task: **"dev: run all"**

## Building

The project uses a single, portable **Node.js script** (`build.js`):

```bash
# Full build
npm run build
# Or: node build.js

# Build steps individually
npm run build:frontend  # Frontend only
npm run build:backend   # Backend only
npm run copy            # Copy frontend to backend only
```

**VS Code**: Press `Ctrl+Shift+B` for default build task

See [BUILD.md](BUILD.md) for detailed build instructions.

## Architecture

- **Frontend**: SvelteKit (static build)
- **Backend**: Go with embedded frontend
- **Single Binary**: The frontend is embedded in the Go binary for easy deployment
- **Data Storage**: Local file-based storage (JSON files) for simplicity (will be replaced with a database in future)
- **API**: RESTful API for managing invoices, clients, and providers

## Project Structure

```
frontend/          # SvelteKit application
backend/           # Go backend with API and embedded UI
  internal/
    api/           # REST API handlers
    ui/            # Embedded frontend (auto-generated)
    storage/       # Data storage layer
build.js           # Single portable build script
```

## API

REST API available at `/api/v1/`:

- `/api/v1/invoices` - Invoice management
- `/api/v1/clients` - Client management
- `/api/v1/providers` - Provider management

## License

GLP-3.0 License. See [LICENSE](LICENSE) for details.
