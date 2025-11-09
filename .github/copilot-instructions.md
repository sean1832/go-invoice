# go-invoice Copilot Instructions

## Architecture Overview

**go-invoice** is a cross-platform invoice management app with a **Go backend** and **SvelteKit frontend**, compiled into a single binary with an embedded UI.

### Key Architecture Decisions

1. **Single Binary Deployment**: Frontend (SvelteKit static build) → Go `embed` package (`//go:embed all:dist`) → `backend/bin/go-invoice[.exe]` (Windows/Linux/macOS)
2. **File-Based Storage**: JSON files in `db/{clients,providers,invoices,email_templates}/` directories (relative to executable or `STORAGE_PATH` env var)
3. **API-First Design**: REST API at `/api/v1/` with ISO date strings (YYYY-MM-DD format); frontend uses relative paths (`/api/v1/*`)
4. **SvelteKit Static Site**: Frontend builds to static HTML (adapter-static) served at root `/` with SPA fallback in `ui.go`
5. **Service Layer Integration**: ChromeDP for PDF generation, SMTP/OAuth2 for email (configured via environment variables)

## Build System

**Single `build.js` Node.js script** orchestrates the entire build (cross-platform):

```bash
npm run build                    # Full: frontend → copy → backend binary
npm run build:frontend           # SvelteKit build
npm run build:backend            # Go binary (output: go-invoice[.exe])
npm run copy                     # Copy frontend/build/ → backend/internal/ui/dist/
```

**VS Code Task** (preferred for development):

- `Ctrl+Shift+B` → "build: full application" (default build task)

**Manual Development Setup**:

```bash
# Terminal 1: Frontend dev server
cd frontend && npm run dev     # http://localhost:5173

# Terminal 2: Backend dev server
cd backend && go run . --dev   # http://localhost:8080 (with CORS for localhost:5173)
```

**Environment Variables** (optional `.env` in `backend/`):

- `STORAGE_PATH`: Override default `db/` location (defaults to `{executable_dir}/db`)
- `DEV_FRONTEND_BASE_URL`: Frontend URL in dev mode (default: `http://localhost:5173`)
- Email config: `SMTP_FROM`, `SMTP_HOST`, `SMTP_PORT`, `SMTP_PASSWORD` (plain auth) or `GOOGLE_OAUTH_CLIENT_ID`, `GOOGLE_OAUTH_CLIENT_SECRET` (OAuth2)

## Critical Developer Workflows

### Local Development

1. Start backend with dev flag: `cd backend && go run . --dev` (enables CORS for localhost:5173)
2. Start frontend: `cd frontend && npm run dev`
3. Frontend at http://localhost:5173 (hot reload), Backend API at http://localhost:8080
4. Frontend sends requests to `/api/v1/*` which proxies/resolves to backend in dev mode

### Production Build & Deployment

1. Run `npm run build` or VS Code task "build: full application"
2. Output: `backend/bin/go-invoice[.exe]` (single binary with embedded UI)
3. **Critical**: Binary looks for data in `db/` relative to its location (or set `STORAGE_PATH` env var)
4. Run: `./backend/bin/go-invoice` → Server starts at http://localhost:8080
5. Storage auto-initializes directories on first run via `storage.NewStorageDir()`

### Adding New Features

**Frontend Route**:

1. Create `.svelte` file in `frontend/src/routes/{path}/+page.svelte`
2. Add to navigation in `frontend/src/routes/+layout.svelte` if needed
3. SPA routing handled automatically by SvelteKit adapter-static

**Backend API Endpoint**:

1. Create handler in `backend/internal/api/handler_resource_{name}.go` or `handler_action_{name}.go`
2. Register in `RegisterRoutesV1()` (e.g., `mux.HandleFunc("/api/v1/resource", h.handleResource)`)
3. Use helper functions: `getResourceByID()`, `createResource()`, `updateResourceByID()`, `deleteResourceByID()`

**Adding PDF/Email Features**:

- PDF: Use `services.NewChromeService()` to render frontend route (e.g., `/invoices/{id}/print`) as PDF via ChromeDP
- Email: Check `h.EmailAuthMethod` (None/Plain/OAuth2), use `services.NewSMTPService()` with attachment support

## Frontend Patterns

### Store Architecture (Single Source of Truth)

- **Centralized stores** in `src/lib/stores/`: All data operations (CRUD) go through Svelte stores
- **Pattern**: Load data → Store updates → UI auto-reacts (no manual state management in components)
- **Key stores**: `invoices`, `clients`, `providers`, `activeProvider` (all with loading states)
- **Services**: Data fetching logic in `src/lib/services/*.service.ts` (uses SvelteKit's `fetch`)
  - Example: `getAllInvoices(fetch)`, `createInvoice(fetch, data)`, `downloadPdf(fetch, id)`
  - All services use `/api/v1/` base path via `apiClient` wrapper
- **Derived stores**: `filteredInvoices` auto-filters based on `invoiceFilters` store
- **Initialization**: Call `loadInvoices()`, `loadClients()`, `loadProviders()` in root layout/page
- See `frontend/src/lib/stores/README.md` for complete store API reference

### Component Structure (Atomic Design)

- **Atoms** (`components/atoms/`): Smallest reusable pieces (date-display, status-badge, currency-display)
- **Molecules** (`components/molecules/`): Combinations of atoms (form groups, cards)
- **Organisms** (`components/organisms/`): Complete features (invoice-form, invoice-display, navbar, shelf, profile-form)
- **UI Components** (`components/ui/`): bits-ui primitives + Tailwind styling (shadcn/ui patterns)
  - Use `$props()` for all component props (never destructure)
  - Use `$bindable()` for two-way binding (e.g., `ref`, `value`, `open`)
  - Import from `bits-ui` (Tabs, Select, Popover, Dialog primitives)

### Date Handling (Critical Pattern)

- **Storage/API format**: ISO YYYY-MM-DD strings (e.g., `"2025-11-03"`)
- **UI format**: `DateValue` from `@internationalized/date` library
- **Sync pattern**: Use `$effect` blocks (NOT utility functions) for two-way binding:
  ```svelte
  <script>
  let dateValue = $state<DateValue | undefined>(undefined);
  $effect(() => {
    if (dateValue && dateValue.toString() !== formString) {
      formString = dateValue.toString();
    }
  });
  </script>
  ```
- **Components**: `DatePicker` auto-closes after selection; `Calendar` for inline date picking
- **CRITICAL RULE**: NEVER import `@internationalized/date` in `.ts` utility files — only in `.svelte` files
  - OK: `frontend/src/lib/helpers/date-helpers.ts` (used by `.svelte` files)
  - NOT OK: Importing in service/store `.ts` files that don't render UI

### Type Safety

- Core types in `src/lib/types/invoice.ts`: `Invoice`, `Party`, `ServiceItem`, `Pricing`, `InvoiceStatus`, `ClientData`, `ProviderData`
- API responses must match these types exactly (snake_case for JSON serialization)
- Date fields: always `string` (ISO format) in types
- Validators in `src/lib/helpers/validators.ts`: `validateInvoice()`, `validateParty()`, `validateLineItem()`

## Backend Patterns

### API Handlers

- Pattern: `Handler struct` with `Context`, `StorageDir`, `FrontendBaseURL`, `EmailAuthMethod` fields
- Register routes in `RegisterRoutesV1()` using path patterns: `/api/v1/resource/{id}`
- JSON serialization with Go's built-in `encoding/json` (snake_case tags)
- Method routing: Check `r.Method` in handler functions (GET/POST/PUT/DELETE)
- Helper functions in `resource_helpers.go`: `getResourceByID()`, `createResource()`, `updateResourceByID()`, `deleteResourceByID()`

### Data Storage

- **Path**: `db/{clients,providers,invoices,email_templates}/` (relative to executable or `STORAGE_PATH` env var)
- **Format**: JSON files with entity IDs as filenames (e.g., `INV-251103.json`, `dingyu_xu.json`)
- **Initialization**: `storage.NewStorageDir(rootDir)` creates directories if missing (auto-setup on first run)
- **Read/Write**: Use `os.ReadFile()` and `os.WriteFile()` with JSON marshaling
- **File operations**: `filepath.Join()` for cross-platform paths
- **Types**: `ClientData` (extends `Party` + `TaxRate`, `EmailTarget`), `ProviderData` (extends `Party` + `PaymentInfo`)

### REST Conventions

- Collection: `GET /api/v1/invoices`, `POST /api/v1/invoices`
- Item: `GET /api/v1/invoices/{id}`, `PUT /api/v1/invoices/{id}`, `DELETE /api/v1/invoices/{id}`
- Special: `GET /api/v1/invoices/count`, `GET /api/v1/invoices/{id}/pdf`, `POST /api/v1/invoices/{id}/email`
- Query params: Support filtering via `?client_id={id}`, `?provider_id={id}`, `?status={status}`, `?date_from={iso}`, `?date_to={iso}`
  - Implemented in `internal/query/query_params.go` and `internal/query/invoice_filters.go`
  - Example: `/api/v1/invoices?status=draft&client_id=dingyu_xu&date_from=2025-01-01`

### Service Integration

- **ChromeService** (`internal/services/chrome.go`): Headless browser for PDF generation

  - Pattern: `NewChromeService()` → `GeneratePDF(url, timeout, paperSize, title)` → `Close()`
  - Renders frontend routes (e.g., `/invoices/{id}/print`) as PDFs via ChromeDP
  - Waits for `#pdf-render-complete` or `#pdf-render-error` elements before generating PDF
  - Paper sizes: `PaperSizeA4`, `PaperSizeA3`, `PaperSizeLetter`

- **SMTPService** (`internal/services/smtp.go`): Email with attachments
  - Pattern: `NewSMTPService(from, host, port, password)` → `SendWithAttachment(...)`
  - Supports plain auth (via `SMTP_PASSWORD`) or OAuth2 (via `GOOGLE_OAUTH_CLIENT_ID/SECRET`)
  - Email templates stored in `db/email_templates/` (e.g., `default.json`)

### CORS & Dev Mode

- CORS middleware in `middleware.go`: `WithCORS(handler, allowedOrigins)`
- Dev mode flag (`--dev`): Enables CORS for `http://localhost:5173` (frontend dev server)
- Production: CORS set to frontend base URL (defaults to `http://localhost:8080`)

## Type System Notes

### Invoice Flow (Frontend → Backend)

1. **Create/Edit**: Invoice editor → form state → `Invoice` type (matches backend `invoice.Invoice`)
2. **Dates**: Issue date & due date stored as ISO strings, synced with `DateValue` UI via `$effect`
3. **Line Items**: Each `ServiceItem` has `date`, `description`, `quantity`, `unit_price`, `total_price`
4. **Pricing**: Auto-calculated `subtotal`, `tax`, `tax_rate`, `total` (matches backend `Pricing`)
5. **Save**: POST/PUT to `/api/v1/invoices` with complete `Invoice` object (snake_case JSON)

### Party (Client/Provider) Flow

1. **Selection**: Dropdown in UI (profile-selector) → `ClientData` or `ProviderData`
2. **Display**: Read-only `Party` object with `name`, `email`, `abn`, `address`, `phone`
3. **Storage**: JSON files in `db/clients/` and `db/providers/` with ID derived from name (e.g., "Dingyu Xu" → `dingyu_xu.json`)
4. **Extra fields**: Clients have `tax_rate`, `email_target`, `email_template_id`; Providers have `payment_info`

### Custom Date Type (Backend)

- `types.Date` wraps `time.Time` with JSON marshaling to/from ISO YYYY-MM-DD strings
- Frontend always sends/receives date strings; backend converts to `types.Date` automatically
- Use `types.Date.String()` for ISO format, `types.Date.Time` for `time.Time` operations

## Important Conventions

- **No `@internationalized/date` imports in utility files** — use directly in `.svelte` files only
- **$effect for runes** — NEVER in `.ts` files, always in `.svelte`
- **Date strings**: Always ISO YYYY-MM-DD format (use `.toISOString().split('T')[0]`)
- **API base**: `/api/v1/` (no trailing slash on collection endpoints)
- **CORS**: Enabled for configured origins in `main.go` (dev mode adds localhost:5173)
- **Port**: Frontend dev 5173, Backend 8080 (hardcoded)

## Key Files Reference

| File                                                                      | Purpose                                            |
| ------------------------------------------------------------------------- | -------------------------------------------------- |
| `build.js`                                                                | Cross-platform build orchestrator                  |
| `backend/main.go`                                                         | Server entry point, route registration             |
| `backend/internal/api/api.go`                                             | API handler interface & route definitions          |
| `backend/internal/storage/storage.go`                                     | File-based storage initialization                  |
| `frontend/src/lib/types/invoice.ts`                                       | Core TypeScript interfaces                         |
| `frontend/src/lib/components/custom/invoice-editor/invoice-editor.svelte` | Invoice creation/editing with date pickers         |
| `frontend/src/lib/components/custom/date-picker/date-picker.svelte`       | Reusable date picker (auto-closes after selection) |

## Testing & Debugging

- **Backend Debug**: Press `F5` → "Launch Backend with Full Build" (builds + debugger) or "Launch Backend (Quick)" (debug only)
- **Frontend Hot Reload**: Active in dev mode; refresh browser after backend changes
- **API Testing**: Use VS Code REST Client or test endpoints at http://localhost:8080/api/v1/
- **Data inspection**: Check `backend/bin/db/` directory for JSON files during development
- **Build verification**: Run `build: full application` task before testing production binary

## Common Pitfalls

1. **Binary finds no data**: Check binary is run from directory containing `db/` folder (or set `STORAGE_PATH` env var)
2. **Date picker closed immediately**: Ensure `Popover.Root bind:open` is implemented
3. **Props not reactive**: Always use `$props()` in component scripts, not destructuring imports
4. **Date sync loops**: Check that sync conditions prevent infinite updates (compare `.toString()`)
5. **ChromeDP timeout**: PDF generation waits for `#pdf-render-complete` element; ensure print pages render this ID
6. **CORS errors**: In dev mode, backend must run with `--dev` flag to allow localhost:5173
