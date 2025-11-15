<div align="center">

<img src=".design/favicon.svg" alt="go-invoice logo" width="120" height="120">

# go-invoice

![GitHub License](https://img.shields.io/github/license/sean1832/go-invoice)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org/)
[![SvelteKit](https://img.shields.io/badge/SvelteKit-5.0-FF3E00?logo=svelte)](https://kit.svelte.dev/)
![GitHub Release](https://img.shields.io/github/v/release/sean1832/go-invoice)

**A simple, no-nonsense invoice management application.**  
Create, manage, and send professional invoices without the bloat of traditional accounting software.

Built with **SvelteKit** and **Go**, packaged as a single executable, No installation, no database to configure. Just run it.

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

### Requirements

**Chrome/Chromium Browser** (required for PDF generation):

### Download & Run

**1. Download the latest release** from [Releases](https://github.com/sean1832/go-invoice/releases/latest) and Extract:

- Windows: `go-invoice-windows-amd64.zip`
- macOS: `go-invoice-macos-amd64.zip`
- macOS (Apple Silicon): `go-invoice-macos-arm64.zip`
- Linux: `go-invoice-linux-amd64.zip`

**2. Download and modify `.env` for configuration** (see [Configuration](#️-configuration)):

Place the `.env` file in the same directory as the executable.

```bash
curl -O https://raw.githubusercontent.com/sean1832/go-invoice/main/backend/.env.example
mv .env.example .env
```

**3. Run the application:**

```bash
# Linux/macOS
chmod +x go-invoice  # Make executable (first time only)
./go-invoice

# Windows (PowerShell)
.\go-invoice.exe
```

**4. Open your browser** to http://localhost:8080

### Quick One-Liner (Linux/macOS)

This command downloads the latest release, extracts it, makes it executable, and sets up a sample `.env` file:

```bash
VERSION=$(curl -s "https://api.github.com/repos/sean1832/go-invoice/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/') && curl -L -o go-invoice.zip "https://github.com/sean1832/go-invoice/releases/download/$VERSION/go-invoice-linux-amd64.zip" && unzip go-invoice.zip && chmod +x go-invoice && curl -O https://raw.githubusercontent.com/sean1832/go-invoice/main/backend/.env.example && mv .env.example .env
```

### Build from Source

```bash
git clone https://github.com/sean1832/go-invoice.git
cd go-invoice/frontend
npm install
cd ..
npm run build
./backend/bin/go-invoice  # or .\backend\bin\go-invoice.exe
```

The server will start at http://localhost:8080

## ⚙️ Configuration

### Basic Configuration

The application works out of the box with no configuration. To enable email sending or customize settings, create a `.env` file in the same directory as the `go-invoice` executable.

**Download the template:**

```bash
curl -O https://raw.githubusercontent.com/sean1832/go-invoice/main/backend/.env.example
mv .env.example .env
```

### Email Setup (Optional)

Choose one of two methods to send invoices via email:

| Feature                 | **App Password**          | **OAuth2** (Recommended)             |
| ----------------------- | ------------------------- | ------------------------------------ |
| **Security**            | Static password in `.env` | Temporary tokens, no password stored |
| **Setup Complexity**    | Simple (5 minutes)        | Moderate (15 minutes)                |
| **Session Persistence** | Permanent                 | 30 days (configurable)               |
| **Best For**            | Quick setup, personal use | Production, better security          |

**📖 Detailed Setup Guides:**

- **[App Password Setup Guide](docs/email-setup-app-password.md)** - Simple SMTP authentication
- **[OAuth2 Setup Guide](docs/email-setup-oauth2.md)** - Secure Google OAuth2 (recommended)

**Quick comparison:**

- **App Password**: Requires 2-Step Verification, creates a static password, stores in `.env`
- **OAuth2**: Requires Google Cloud project, OAuth consent screen setup, temporary access tokens

> [!NOTE]
> Both methods require setting up credentials in your Google Account or Google Cloud Console. See the detailed guides above for step-by-step instructions.

### Storage Location

By default, data is stored in `db/` **relative to the executable**.

```bash
# Custom storage location
export STORAGE_PATH=/path/to/your/data  # Linux/macOS
$env:STORAGE_PATH="C:\path\to\data"     # Windows PowerShell
```

> [!WARNING]
> If you move the binary or run it from a different directory, your data will appear missing unless you set an absolute `STORAGE_PATH`. Always run the binary from its directory or use absolute paths.

### Environment Variables

See [`.env.example`](backend/.env.example) for all available configuration options including:

- `PORT` - Server port (default: `8080`)
- `PUBLIC_URL` - Public-facing URL for OAuth2 callbacks (default: `http://localhost:{PORT}`)
- `SESSION_SECRET` - **Critical for production** (see below)
- `SESSION_MAX_AGE` - Session duration in seconds (default: 2592000 = 30 days)
- `IS_PROD` - Enable production mode (default: `false`)

> [!IMPORTANT]
> **For Production:** Set a persistent `SESSION_SECRET` to prevent users from being logged out when the server restarts:
>
> ```bash
> # Generate a secure secret
> openssl rand -base64 32
> ```
>
> Add to `.env`:
>
> ```env
> SESSION_SECRET=your-generated-secret-here
> ```
>
> Without this, a new secret is generated on each restart and all previous sessions become invalid.

### Command-Line Options

```bash
./go-invoice --port 3000        # Custom port
./go-invoice --dev              # Development mode (enables CORS)
./go-invoice --db /custom/path  # Custom database path
```

## 🚀 Production Deployment

### Pre-Deployment Checklist

Before deploying to production, ensure you've configured:

- [ ] **SESSION_SECRET** - Generate with `openssl rand -base64 32` and add to `.env`
- [ ] **PUBLIC_URL** - Set to your HTTPS domain (e.g., `https://yourdomain.com`)
- [ ] **IS_PROD** - Set to `true` for secure cookies
- [ ] **Chrome/Chromium** - Verify installation with `google-chrome --version` or `chromium --version`
- [ ] **Firewall/Reverse Proxy** - Configure port 8080 access or set up reverse proxy
- [ ] **Email Authentication** - Complete either [App Password](docs/email-setup-app-password.md) or [OAuth2](docs/email-setup-oauth2.md) setup
- [ ] **STORAGE_PATH** - Set absolute path for data storage (optional)

### Session Management

Sessions persist for 30 days by default (`SESSION_MAX_AGE=2592000` seconds).

**Why `SESSION_SECRET` matters:**

- If not set, a random secret is generated on each server restart
- All users are logged out when the server restarts
- Sessions become invalid between deployments

**Generate a persistent secret:**

```bash
openssl rand -base64 32
```

Add to `.env`:

```env
SESSION_SECRET=AbCdEf123456...  # Your generated secret
```

### Environment Variables for Production

```env
# Server Configuration
PORT=8080
PUBLIC_URL=https://yourdomain.com
IS_PROD=true

# Session (CRITICAL - generate with: openssl rand -base64 32)
SESSION_SECRET=your-persistent-secret-here
SESSION_MAX_AGE=2592000

# Email - Choose ONE method
# Method 1: OAuth2 (Recommended)
GOOGLE_OAUTH_CLIENT_ID=your-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-client-secret
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# Method 2: App Password (Not recommended for production)
# SMTP_FROM=your-email@gmail.com
# SMTP_HOST=smtp.gmail.com
# SMTP_PORT=587
# SMTP_PASSWORD=your-app-password

# Storage
STORAGE_PATH=/var/lib/go-invoice/db
```

## 📝 License

GPL-3.0 License. See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Built with:

- [SvelteKit](https://kit.svelte.dev/) - Web framework
- [Go](https://golang.org/) - Backend language
- [ChromeDP](https://github.com/chromedp/chromedp) - PDF generation
- [Goth](https://github.com/markbates/goth) - OAuth2 provider
- [shadcn-svelte](https://shadcn-svelte.com/) - UI components
