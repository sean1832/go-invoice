<div align="center">

<img src="design/favicon.svg" alt="go-invoice logo" width="120" height="120">

# go-invoice

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org/)
[![SvelteKit](https://img.shields.io/badge/SvelteKit-5.0-FF3E00?logo=svelte)](https://kit.svelte.dev/)
[![Release](https://img.shields.io/github/v/release/sean1832/go-invoice)](https://github.com/sean1832/go-invoice/releases/latest)
[![Build Status](https://img.shields.io/github/actions/workflow/status/sean1832/go-invoice/build.yml?branch=main)](https://github.com/sean1832/go-invoice/actions)

**A simple, no-nonsense invoice management application.**  
Create, manage, and send professional invoices without the bloat of traditional accounting software.

Built with **SvelteKit** and **Go**, packaged as a single executable-no installation, no dependencies, no database to configure. Just run it.

[Features](#-features) • [Quick Start](#-quick-start) • [Configuration](#️-configuration) • [API](#-api) • [Development](#️-development)

</div>

---

> [!WARNING]
> This project is currently in development. Some features may be incomplete or experimental.

## ✨ Features

- 📝 **Create & Manage Invoices** - Simple forms, automatic calculations, professional layouts
- 👥 **Client & Provider Management** - Store contact details, payment info, and preferences
- 📄 **PDF Generation** - Export invoices as PDFs with one click (powered by headless Chrome)
- 📧 **Email Integration** - Send invoices directly to clients via SMTP or Gmail OAuth2
- 🗂️ **File-Based Storage** - No database setup required everything stored as JSON files
- 🔌 **REST API** - Integrate with your existing tools and workflows
- 🚀 **Single Binary** - Frontend and backend bundled together for easy deployment
- 🌐 **Cross-Platform** - Works on Windows, macOS, and Linux

## 🚀 Quick Start

### Download & Run

**Desktop/Local:**

1. **Download** the latest release for your platform from [Releases](https://github.com/sean1832/go-invoice/releases/latest)

   - Windows: `go-invoice.exe`
   - macOS/Linux: `go-invoice`

2. **Run** the application:

   ```bash
   # Windows
   .\go-invoice.exe

   # macOS/Linux
   ./go-invoice
   ```

3. **Open** your browser to http://localhost:8080

That's it! Your data will be stored in a `db/` folder next to the executable.

**Headless Server (Linux):**

```bash
# Download the latest release
VERSION=$(curl -s "https://api.github.com/repos/sean1832/go-invoice/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
curl -L -o go-invoice "https://github.com/sean1832/go-invoice/releases/download/$VERSION/go-invoice"
chmod +x go-invoice

# Run the application
./go-invoice
```

> [!IMPORTANT] Chrome/Chromium Required for PDF Generation
>
> The server must have Chrome or Chromium installed for PDF generation to work. Install it with:
>
> ```bash
> # Ubuntu/Debian
> sudo apt-get update && sudo apt-get install -y chromium-browser
>
> # CentOS/RHEL
> sudo yum install -y chromium
>
> # Or use Google Chrome
> wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
> sudo dpkg -i google-chrome-stable_current_amd64.deb
> ```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/sean1832/go-invoice.git
cd go-invoice

# Install dependencies
npm install

# Build the application
npm run build

# Run it
./backend/bin/go-invoice      # Linux/Mac
.\backend\bin\go-invoice.exe  # Windows
```

The server will start at http://localhost:8080

## ⚙️ Configuration

### Storage Location

By default, data is stored in `db/` relative to the executable. To use a custom location:

```bash
export STORAGE_PATH=/path/to/your/data  # Linux/Mac
set STORAGE_PATH=C:\path\to\your\data   # Windows

./go-invoice
```

### Email Setup (Optional)

To send invoices via email, configure one of the following methods:

**Using `.env` file** (Recommended)

Copy the example configuration file and edit it:

```bash
cp .env.example .env
# Edit .env with your preferred editor
nano .env  # or vim, code, etc.
```

Example `.env` configuration:

```env
# SMTP Configuration (Option 1: Simple but less secure)
SMTP_FROM=your-email@example.com
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_PASSWORD=your-app-password

# OR Google OAuth2 (Option 2: Recommended for Gmail)
# GOOGLE_OAUTH_CLIENT_ID=your-client-id
# GOOGLE_OAUTH_CLIENT_SECRET=your-client-secret
```

**Using Environment Variables**

**Option 1: SMTP with Password** (Simple but less secure)

```bash
export SMTP_FROM="your-email@example.com"
export SMTP_HOST="smtp.gmail.com"
export SMTP_PORT="587"
export SMTP_PASSWORD="your-app-password"
```

**Option 2: Gmail OAuth2** (Recommended)

```bash
export GOOGLE_OAUTH_CLIENT_ID="your-client-id"
export GOOGLE_OAUTH_CLIENT_SECRET="your-client-secret"
```

### Advanced Options

- `DEV_FRONTEND_BASE_URL` - Frontend URL for development mode (default: `http://localhost:5173`)
- `--dev` flag - Enable development mode with CORS for local frontend
- `--port` flag - Custom port (default: `8080`)

```bash
./go-invoice --port 3000 --dev
```

## 🏗️ Architecture

go-invoice is built with simplicity in mind:

- **Frontend**: SvelteKit static site (modern, reactive UI)
- **Backend**: Go HTTP server with embedded frontend
- **Storage**: JSON files (no database required)
- **PDF Engine**: ChromeDP (headless Chrome for pixel-perfect PDFs)
- **Deployment**: Single binary—frontend and backend bundled together

**Why this stack?**

- ✅ No database to configure or maintain
- ✅ No separate frontend deployment needed
- ✅ Easy to backup (just copy the `db/` folder)
- ✅ Simple to version control your data
- ✅ Perfect for freelancers and small businesses

## 📡 API

RESTful API available at `/api/v1/` for integration with other tools:

- `GET/POST /api/v1/invoices` - List and create invoices
- `GET/PUT/DELETE /api/v1/invoices/{id}` - Manage individual invoices
- `GET /api/v1/invoices/{id}/pdf` - Generate PDF for an invoice
- `POST /api/v1/invoices/{id}/email` - Send invoice via email
- `GET/POST/DELETE /api/v1/clients` - Manage clients
- `GET/POST/DELETE /api/v1/providers` - Manage service providers


## 🛠️ Development

Want to contribute or customize the app?

```bash
# Install dependencies
npm install

# Run in development mode
npm run dev:frontend  # Terminal 1 - http://localhost:5173
npm run dev:backend   # Terminal 2 - http://localhost:8080 (with --dev flag)
```

Or use the VS Code task: **"dev: run all"**

**Build Commands:**

```bash
npm run build           # Full build (frontend + backend)
npm run build:frontend  # Frontend only
npm run build:backend   # Backend only
```

See [BUILD.md](BUILD.md) for detailed development instructions.

## 📦 Project Structure

```
go-invoice/
├── frontend/              # SvelteKit application
│   ├── src/
│   │   ├── routes/        # Pages and layouts
│   │   └── lib/           # Components, stores, services
│   └── build/             # Static build output
├── backend/               # Go backend
│   ├── main.go            # Entry point
│   ├── internal/
│   │   ├── api/           # REST API handlers
│   │   ├── services/      # PDF & email services
│   │   ├── storage/       # File-based storage layer
│   │   └── ui/            # Embedded frontend (auto-generated)
│   └── bin/
│       ├── go-invoice     # Compiled binary
│       └── db/            # Data storage (JSON files)
└── build.js               # Cross-platform build script
```

## 📝 License

GPL-3.0 License. See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Built with:

- [SvelteKit](https://kit.svelte.dev/) - Web framework
- [Go](https://golang.org/) - Backend language
- [ChromeDP](https://github.com/chromedp/chromedp) - PDF generation
- [shadcn-svelte](https://shadcn-svelte.com/) - UI components
