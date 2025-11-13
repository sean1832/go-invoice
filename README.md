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

- **Ubuntu/Debian**:

  ```bash
  sudo apt-get update && sudo apt-get install -y chromium-browser
  ```

- **Windows** (using Chocolatey):

  ```powershell
  choco install googlechrome
  ```

  Or download [Google Chrome](https://www.google.com/chrome/) manually

- **macOS** (using Homebrew):
  ```bash
  brew install --cask google-chrome
  ```

### Download & Run

**1. Download the latest release** from [Releases](https://github.com/sean1832/go-invoice/releases/latest):

- Windows: `go-invoice-windows-amd64.zip`
- macOS: `go-invoice-macos-amd64.zip`
- macOS (Apple Silicon): `go-invoice-macos-arm64.zip`
- Linux: `go-invoice-linux-amd64.zip`

**2. Extract the archive** to a folder of your choice

**3. Run the application:**

```bash
# Linux/macOS
chmod +x go-invoice  # Make executable (first time only)
./go-invoice

# Windows (PowerShell)
.\go-invoice.exe
```

**4. Open your browser** to http://localhost:8080

### Directory Structure

After extraction, your directory should look like this:

```
your-app-folder/
├── go-invoice          # The executable (go-invoice.exe on Windows)
├── .env                # Configuration file (create this for email/custom settings)
└── db/                 # Data storage (auto-created on first run)
    ├── clients/
    ├── invoices/
    ├── providers/
    └── email_templates/
```

> [!IMPORTANT]
> **The `.env` file must be in the same directory as the executable.** Data is stored in the `db/` folder next to the binary unless you set a custom `STORAGE_PATH` or use `--db <path>` flag.

### Quick One-Liner (Linux/macOS)

```bash
VERSION=$(curl -s "https://api.github.com/repos/sean1832/go-invoice/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/') && curl -L -o go-invoice.zip "https://github.com/sean1832/go-invoice/releases/download/$VERSION/go-invoice-linux-amd64.zip" && unzip go-invoice.zip && chmod +x go-invoice && ./go-invoice
```

### Build from Source

See [BUILD.md](BUILD.md) for detailed instructions.

```bash
git clone https://github.com/sean1832/go-invoice.git
cd go-invoice
npm install
npm run build
./backend/bin/go-invoice  # or .\backend\bin\go-invoice.exe on Windows
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

### Data Backup

All data is stored in the `db/` folder. To backup:

```bash
# Backup
tar -czf invoice-backup-$(date +%Y%m%d).tar.gz db/

# Restore
tar -xzf invoice-backup-20251114.tar.gz
```

Or simply copy the `db/` folder to another location.

### Reverse Proxy Configuration

**Nginx Example:**

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

For HTTPS (recommended):

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**Caddy (automatic HTTPS):**

```caddy
yourdomain.com {
    reverse_proxy localhost:8080
}
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

## 🔧 Troubleshooting

### PDF Generation Fails

**Symptom:** Error when trying to download or email invoice PDFs

**Solutions:**

1. Verify Chrome/Chromium is installed:

   ```bash
   google-chrome --version    # or
   chromium --version
   ```

2. Install if missing (see [Requirements](#requirements))

3. Check application logs for "could not launch browser" errors

4. On Linux servers, ensure required dependencies are installed:
   ```bash
   # Ubuntu/Debian
   sudo apt-get install -y fonts-liberation libnss3 libatk-bridge2.0-0 libx11-xcb1
   ```

### Logged Out After Every Restart

**Symptom:** Users must re-authenticate after server restarts

**Cause:** `SESSION_SECRET` not set or changes between restarts

**Solution:**

1. Generate a persistent secret:

   ```bash
   openssl rand -base64 32
   ```

2. Add to `.env`:

   ```env
   SESSION_SECRET=your-generated-secret-here
   ```

3. Restart the application

4. Verify the `.env` file is in the same directory as the executable

### OAuth "redirect_uri_mismatch" Error

**Symptom:** Error during Google sign-in: "redirect_uri_mismatch"

**Cause:** `PUBLIC_URL` doesn't match the redirect URI configured in Google Cloud Console

**Solution:**

1. Check your `PUBLIC_URL` in `.env`:

   ```env
   PUBLIC_URL=http://localhost:8080  # or your production URL
   ```

2. The redirect URI should be:

   ```
   {PUBLIC_URL}/api/v1/mailer/auth/google/callback
   ```

3. Verify in Google Cloud Console:

   - Go to [APIs & Services > Credentials](https://console.cloud.google.com/apis/credentials)
   - Click your OAuth client
   - Check "Authorized redirect URIs" matches exactly (including `http`/`https` and port)

4. Common mistakes:
   - Missing port number (`:8080`)
   - Wrong protocol (`http` vs `https`)
   - Trailing slash in `PUBLIC_URL`
   - Missing `/api/v1/mailer/auth/google/callback` path

### Email Sending Fails

**For App Password method:**

1. Verify 2-Step Verification is enabled on your Google Account
2. Check the App Password has no spaces in `.env`
3. Confirm `SMTP_FROM` matches your Gmail address
4. Ensure `SMTP_HOST=smtp.gmail.com` and `SMTP_PORT=587`

**For OAuth2 method:**

1. Verify you've authenticated in Settings
2. Check `GOOGLE_OAUTH_CLIENT_ID` and `GOOGLE_OAUTH_CLIENT_SECRET` are correct
3. Ensure `PUBLIC_URL` matches Google Cloud Console redirect URI
4. Do NOT set `SMTP_PASSWORD` (conflicts with OAuth2)
5. Try signing out and signing in again in Settings

### Data Not Found After Moving Binary

**Symptom:** Invoices, clients, or providers missing after moving the executable

**Cause:** Data is stored in `db/` relative to the executable's location

**Solutions:**

1. **Option 1:** Run the binary from its original directory

   ```bash
   cd /path/to/original/location
   ./go-invoice
   ```

2. **Option 2:** Set absolute `STORAGE_PATH`:

   ```bash
   export STORAGE_PATH=/absolute/path/to/db  # Linux/macOS
   $env:STORAGE_PATH="C:\absolute\path\to\db"  # Windows PowerShell
   ```

3. **Option 3:** Move the `db/` folder to the new location

### Port 8080 Already in Use

**Symptom:** "address already in use" error on startup

**Solution:**

1. Change the port in `.env`:

   ```env
   PORT=3000  # or any available port
   ```

2. Or use command-line flag:

   ```bash
   ./go-invoice --port 3000
   ```

3. Update `PUBLIC_URL` if using OAuth2:
   ```env
   PUBLIC_URL=http://localhost:3000
   ```

### Application Won't Start

1. Check file permissions (executable flag):

   ```bash
   chmod +x go-invoice
   ```

2. Verify `.env` syntax (no quotes around values unless they contain spaces)

3. Check application logs for error messages

4. Ensure no conflicting environment variables are set

### Still Having Issues?

- Check the [detailed email setup guides](docs/)
- Review the [`.env.example`](backend/.env.example) for configuration reference
- Open an issue on [GitHub Issues](https://github.com/sean1832/go-invoice/issues)

## 📡 API

RESTful API available at `/api/v1/`:

- `GET/POST /api/v1/invoices` - List and create invoices
- `GET/PUT/DELETE /api/v1/invoices/{id}` - Manage individual invoices
- `GET /api/v1/invoices/{id}/pdf` - Generate PDF
- `POST /api/v1/invoices/{id}/email` - Send invoice via email
- `GET/POST/DELETE /api/v1/clients` - Manage clients
- `GET/POST/DELETE /api/v1/providers` - Manage providers

## 🏗️ Architecture

- **Frontend**: SvelteKit static site
- **Backend**: Go HTTP server with embedded frontend
- **Storage**: JSON files (no database required)
- **PDF Engine**: ChromeDP (headless Chrome)
- **Deployment**: Single binary with embedded UI

## 🛠️ Development

See [BUILD.md](BUILD.md) for detailed development and build instructions.

**Quick Start:**

```bash
npm install                 # Install dependencies
npm run dev:frontend        # Terminal 1 - http://localhost:5173
npm run dev:backend         # Terminal 2 - http://localhost:8080
npm run build               # Full build (frontend + backend)
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
