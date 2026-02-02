<div align="center">

<img src=".design/favicon.svg" alt="go-invoice logo" width="120" height="120">

# go-invoice

![GitHub License](https://img.shields.io/github/license/sean1832/go-invoice)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org/)
[![SvelteKit](https://img.shields.io/badge/SvelteKit-5.0-FF3E00?logo=svelte)](https://kit.svelte.dev/)
![GitHub Release](https://img.shields.io/github/v/release/sean1832/go-invoice)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker)](https://hub.docker.com/r/sean1832/go-invoice)

**A simple, no-nonsense invoice management application.**  
Create, manage, and send professional invoices without the bloat of traditional accounting software.

[Quick Start](#-quick-start) • [Configuration](#️-configuration) • [Email Setup](#-email-setup) • [Development](#️-development)

</div>

---

> [!WARNING]
> This project is currently in development. Some features may be incomplete or experimental.

## ✨ Features

- 📝 **Create & Manage Invoices** - Simple forms, automatic calculations, professional layouts
- 👥 **Client & Provider Management** - Store contact details, payment info, and preferences
- 📄 **PDF Generation** - Export invoices as PDFs with one click
- 📧 **Email Integration** - Send invoices directly to clients via SMTP or Gmail OAuth2
- 🗂️ **File-Based Storage** - No database setup required—everything stored as JSON files
- 🔌 **REST API** - Integrate with your existing tools and workflows
- � **Docker Ready** - One command deployment with Docker Compose

## 🚀 Quick Start

Get up and running in under a minute with Docker.

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/)

### 1. Create Docker Compose File

Create a `docker-compose.yml`:

```yaml
services:
  app:
    image: sean1832/go-invoice:latest
    container_name: go-invoice-app
    environment:
      - CHROME_REMOTE_URL=http://chrome:9222
      - CHROME_RENDER_URL=http://app:8080
      - STORAGE_PATH=/data
      - PORT=8080
    ports:
      - "8080:8080"
    volumes:
      - ./db:/data
    networks:
      - invoice-network
    depends_on:
      - chrome

  chrome:
    image: gcr.io/zenika-hub/alpine-chrome:latest
    container_name: go-invoice-chrome
    user: root
    entrypoint:
      - /bin/sh
      - -c
      - |
        apk add --no-cache font-noto font-noto-cjk ttf-dejavu font-noto-emoji terminus-font
        fc-cache -f
        chromium-browser \
        --no-sandbox \
        --disable-dev-shm-usage \
        --disable-gpu \
        --headless \
        --remote-debugging-address=0.0.0.0 \
        --remote-debugging-port=9222 \
        --remote-allow-origins=*
    ports:
      - "9222:9222"
    shm_size: "2gb"
    networks:
      - invoice-network

networks:
  invoice-network:
    driver: bridge
```

### 2. Start the Application

```bash
docker compose up -d
```

### 3. Open Your Browser

Navigate to **http://localhost:8080** - That's it! 🎉

Your data will be stored in `./db` on your host machine.

---

## ⚙️ Configuration

### Environment Variables

Add environment variables to the `app` service in your `docker-compose.yml`:

```yaml
environment:
  # Required
  - CHROME_REMOTE_URL=http://chrome:9222
  - CHROME_RENDER_URL=http://app:8080
  - STORAGE_PATH=/data
  - PORT=8080

  # Production (recommended)
  - PUBLIC_URL=https://yourdomain.com
  - IS_PROD=true
  - SESSION_SECRET=your-secure-secret  # Generate with: openssl rand -base64 32

  # Optional
  - SESSION_MAX_AGE=2592000  # 30 days in seconds
```

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `PUBLIC_URL` | Public URL for OAuth callbacks | `http://localhost:{PORT}` |
| `SESSION_SECRET` | Secret for session encryption | Auto-generated (not persistent) |
| `SESSION_MAX_AGE` | Session duration in seconds | `2592000` (30 days) |
| `IS_PROD` | Enable production mode (secure cookies) | `false` |
| `STORAGE_PATH` | Data storage path inside container | `/data` |

> [!IMPORTANT]
> **For Production:** Always set `SESSION_SECRET` to a persistent value. Without it, all users are logged out when the container restarts.
> ```bash
> openssl rand -base64 32
> ```

---

## 📧 Email Setup

Choose one of two methods to send invoices via email:

| Feature | **App Password** | **OAuth2** (Recommended) |
|---------|------------------|--------------------------|
| **Security** | Static password | Temporary tokens |
| **Setup Time** | ~5 minutes | ~15 minutes |
| **Best For** | Personal use | Production |

### Option 1: App Password (Simple)

```yaml
environment:
  - SMTP_FROM=your-email@gmail.com
  - SMTP_PASSWORD=your-app-password
  - SMTP_HOST=smtp.gmail.com
  - SMTP_PORT=587
```

📖 [App Password Setup Guide](docs/email-setup-app-password.md)

### Option 2: OAuth2 (Recommended)

```yaml
environment:
  - GOOGLE_OAUTH_CLIENT_ID=your-client-id.apps.googleusercontent.com
  - GOOGLE_OAUTH_CLIENT_SECRET=your-client-secret
  - SMTP_HOST=smtp.gmail.com
  - SMTP_PORT=587
```

📖 [OAuth2 Setup Guide](docs/email-setup-oauth2.md)

---

## 🚀 Production Deployment

### Pre-Deployment Checklist

- [ ] Set `SESSION_SECRET` with a generated secret
- [ ] Set `PUBLIC_URL` to your HTTPS domain
- [ ] Set `IS_PROD=true` for secure cookies
- [ ] Configure email authentication (App Password or OAuth2)
- [ ] Set up a reverse proxy (Nginx/Caddy) for SSL termination

### Production Docker Compose

```yaml
services:
  app:
    image: sean1832/go-invoice:latest
    container_name: go-invoice-app
    environment:
      - CHROME_REMOTE_URL=http://chrome:9222
      - CHROME_RENDER_URL=http://app:8080
      - STORAGE_PATH=/data
      - PORT=8080
      - PUBLIC_URL=https://yourdomain.com
      - IS_PROD=true
      - SESSION_SECRET=your-generated-secret-here
      # Email (OAuth2)
      - GOOGLE_OAUTH_CLIENT_ID=your-client-id
      - GOOGLE_OAUTH_CLIENT_SECRET=your-client-secret
      - SMTP_HOST=smtp.gmail.com
      - SMTP_PORT=587
    ports:
      - "8080:8080"
    volumes:
      - ./db:/data
    networks:
      - invoice-network
    depends_on:
      - chrome
    restart: unless-stopped

  chrome:
    image: gcr.io/zenika-hub/alpine-chrome:latest
    container_name: go-invoice-chrome
    user: root
    entrypoint:
      - /bin/sh
      - -c
      - |
        apk add --no-cache font-noto font-noto-cjk ttf-dejavu font-noto-emoji terminus-font
        fc-cache -f
        chromium-browser \
        --no-sandbox \
        --disable-dev-shm-usage \
        --disable-gpu \
        --headless \
        --remote-debugging-address=0.0.0.0 \
        --remote-debugging-port=9222 \
        --remote-allow-origins=*
    shm_size: "2gb"
    networks:
      - invoice-network
    restart: unless-stopped

networks:
  invoice-network:
    driver: bridge
```

### Reverse Proxy (Nginx)

For SSL termination, add this to your Nginx config:

```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # For large PDFs
        client_max_body_size 0;
        proxy_request_buffering off;
        proxy_buffering off;
    }
}
```

---

## 🛠️ Development

### Running from Source

```bash
git clone https://github.com/sean1832/go-invoice.git
cd go-invoice

# Install frontend dependencies
cd frontend && npm install && cd ..

# Build and run
npm run build
./backend/bin/go-invoice  # or go-invoice.exe on Windows
```

### Development Mode

```bash
# Terminal 1: Backend
cd backend && go run . --dev

# Terminal 2: Frontend (hot reload)
cd frontend && npm run dev
```

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

### Building from Dockerfile

```bash
docker build -t go-invoice .
```

---

## 📝 License

GPL-3.0 License. See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Built with [SvelteKit](https://kit.svelte.dev/), [Go](https://golang.org/), [ChromeDP](https://github.com/chromedp/chromedp), [Goth](https://github.com/markbates/goth), and [shadcn-svelte](https://shadcn-svelte.com/).
