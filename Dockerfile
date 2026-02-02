# Build Stage 1: Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app

# Install dependencies (cache optimized)
COPY frontend/package*.json ./
RUN npm ci

# Build frontend
COPY frontend/ ./
RUN npm run build

# Build Stage 2: Backend
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app

# Install dependencies (cache optimized)
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./

# Copy built frontend assets to backend location
COPY --from=frontend-builder /app/build ./internal/ui/dist

# Build Go binary
# CGO_ENABLED=0 is critical for Alpine compatibility
# -s -w strips debug info for smaller size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o go-invoice .

# Build Stage 3: Final Image
FROM alpine:latest
WORKDIR /app

# Install basic runtime dependencies (if needed) and ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=backend-builder /app/go-invoice .

# Create data directory structure
RUN mkdir -p db/clients db/providers db/invoices db/email_templates

# Set production environment variables
ENV STORAGE_PATH=./db
ENV PORT=8080

EXPOSE 8080

CMD ["./go-invoice"]
