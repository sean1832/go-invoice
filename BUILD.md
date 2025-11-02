# Build & Development Guide

## Prerequisites

- Go 1.25.3 or later
- Node.js 16+ (for frontend and build scripts)
- VS Code (recommended)

## Quick Build

The project uses a **single portable Node.js script** (`build.js`) that works on Windows, Linux, and macOS:

```bash
# Full build (frontend + backend)
npm run build
# Or: node build.js

# Individual build steps
npm run build:frontend  # node build.js frontend
npm run build:backend   # node build.js backend
npm run copy            # node build.js copy
```

Output: `backend/bin/invoice-app` (or `.exe` on Windows)

## VS Code Tasks (Cross-Platform)

This project uses VS Code tasks for all build operations, making it cross-platform and developer-friendly.

### Available Tasks

Run tasks via: **Terminal → Run Task** or `Ctrl+Shift+P` → "Tasks: Run Task"

#### Individual Tasks

- **`frontend: install dependencies`** - Install npm packages
- **`frontend: build`** - Build the SvelteKit frontend (static output)
- **`frontend: dev`** - Start frontend dev server (http://localhost:5173)
- **`copy: frontend to backend`** - Copy build to backend for embedding
- **`go: build backend`** - Build the Go binary with embedded UI
- **`go: run backend`** - Run backend server (http://localhost:8080)

#### Composite Tasks

- **`build: full application`** ⭐ _(Default Build)_ - Complete production build

  - Builds frontend
  - Copies to backend
  - Builds Go binary with embedded UI
  - Output: `backend/bin/invoice-app` (or `.exe` on Windows)

- **`dev: run all`** - Run both dev servers in parallel
  - Frontend at http://localhost:5173
  - Backend at http://localhost:8080

### Quick Start

**Production Build:**

```bash
# Via npm (recommended)
npm run build

# Via VS Code: Ctrl+Shift+B (default build task)

# Or directly
node build.js
```

**Development Mode:**

```bash
# Run task: "dev: run all" in VS Code
# Or manually:
npm run dev:frontend  # Terminal 1
npm run dev:backend   # Terminal 2
```

**Run Production Binary:**

```bash
./backend/bin/invoice-app      # Linux/Mac
.\backend\bin\invoice-app.exe  # Windows
```

## Manual Build (Alternative)

### Frontend Only

```bash
cd frontend
npm install
npm run build
```

### Copy Build to Backend

```bash
# From project root
npm run copy
# Or: node build.js copy
```

### Backend Only

```bash
cd backend
go build -o bin/invoice-app .
```

### Full Build (Step by Step)

```bash
# 1. Build frontend
npm run build:frontend
# Or: node build.js frontend

# 2. Copy to backend
npm run copy
# Or: node build.js copy

# 3. Build backend
npm run build:backend
# Or: node build.js backend
```

## Architecture

The frontend is built as a static site and embedded into the Go binary using Go's `embed` package:

```
frontend/build/          → SvelteKit static output
backend/internal/ui/dist/ → Embedded in Go binary
backend/bin/invoice-app   → Single binary with UI + API
```

## Debugging

Use VS Code's debug configurations:

- **Launch Backend** - Builds and runs with debugger
- **Debug Backend (no rebuild)** - Fast debugging without rebuild

Press `F5` or go to **Run → Start Debugging**

## Environment

The embedded UI is served at the root path (`/`), while the API is available at `/api/v1/`:

- `http://localhost:8080/` - Web UI (embedded)
- `http://localhost:8080/api/v1/` - REST API
