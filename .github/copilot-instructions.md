# go-invoice Copilot Instructions

## Architecture Overview

**go-invoice** is a cross-platform invoice management app with a **Go backend** and **SvelteKit frontend**, compiled into a single binary with an embedded UI.

### Key Architecture Decisions

1. **Single Binary Deployment**: Frontend (SvelteKit static build) → Go `embed` package → `backend/bin/invoice-app` (Windows/Linux/macOS)
2. **File-Based Storage**: JSON files in `db/{clients,providers,invoices,smtp}` directories (relative to executable)
3. **API-First Design**: REST API at `/api/v1/` with ISO date strings (YYYY-MM-DD format)
4. **Svelte 5 Runes**: Uses `$state`, `$derived`, `$effect` directly (no imports required in `.svelte` files)
5. **Component Library**: bits-ui + custom Tailwind components with shadcn/ui patterns

## Build System

**Single `build.js` Node.js script** orchestrates the entire build (cross-platform):

```bash
npm run build                    # Full: frontend → copy → backend binary
npm run build:frontend           # SvelteKit build
npm run build:backend            # Go binary
npm run copy                     # Copy frontend/build/ → backend/internal/ui/dist/
```

**VS Code Tasks** (preferred for development):

- `Ctrl+Shift+B` → Default build task
- `dev: run all` → Parallel frontend (5173) + backend (8080) dev servers

## Critical Developer Workflows

### Local Development

1. Run `dev: run all` task to start both servers
2. Frontend at http://localhost:5173 (with hot reload)
3. Backend API at http://localhost:8080
4. Frontend connects to backend via `/api/v1/` endpoints

### Production Build

1. Run `build: full application` task (VS Code) or `npm run build`
2. Single binary at `backend/bin/invoice-app`
3. **Important**: Binary looks for data in `db/` relative to executable location

### Adding Routes

- **Frontend**: Create `.svelte` files under `frontend/src/routes/`
- **Backend**: Add handlers to `backend/internal/api/handler_resource_*.go` and register in `RegisterRoutesV1()`

## Frontend Patterns

### Store Architecture (Single Source of Truth)

- **Centralized stores** in `src/lib/stores/`: All data operations (CRUD) go through Svelte stores
- **Pattern**: Load data → Store updates → UI auto-reacts (no manual state management in components)
- **Key stores**: `invoices`, `clients`, `providers`, `activeProvider` (all with loading states)
- **API-ready**: Stores use localStorage now; swap with fetch calls to `/api/v1/` endpoints (same signatures)
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

### Date Handling (Svelte 5)

- **String format**: ISO YYYY-MM-DD (e.g., `"2025-11-03"`) for API/storage
- **UI format**: Use `DateValue` from `@internationalized/date` library
- **Sync pattern**: Use `$effect` blocks (NOT utility functions) for two-way binding:
  ```svelte
  let dateValue = $state<DateValue | undefined>(undefined);
  $effect(() => {
    if (dateValue && dateValue.toString() !== formString) {
      formString = dateValue.toString();
    }
  });
  ```
- **Components**: `DatePicker` auto-closes after selection; `Calendar` for custom date picking
- **CRITICAL**: NEVER import `@internationalized/date` in `.ts` utility files — only in `.svelte` files

### Type Safety

- Core types in `src/lib/types/invoice.ts`: `Invoice`, `Party`, `ServiceItem`, `Pricing`, `InvoiceStatus`
- API responses must match these types
- Date fields: always `string` (ISO format)
- Validators in `src/lib/utils/validators.ts`: `validateInvoice()`, `validateParty()`, `validateLineItem()`

## Backend Patterns

### API Handlers

- Pattern: `Handler struct` with `Context` and `StorageDir` fields
- Register routes in `RegisterRoutesV1()` using path patterns: `/api/v1/resource/{id}`
- JSON serialization with Go's built-in `encoding/json`
- Method routing: Check `r.Method` in handler functions (GET/POST/PUT/DELETE)

### Data Storage

- **Path**: `db/{clients,providers,invoices,smtp}/` (relative to executable)
- **Format**: JSON files with entity IDs as filenames (e.g., `INV-251103.json`)
- **Initialization**: `storage.NewStorageDir()` creates directories if missing (auto-setup on first run)
- **Read/Write**: Use `os.ReadFile()` and `os.WriteFile()` with JSON marshaling
- **File operations**: `filepath.Join()` for cross-platform paths

### REST Conventions

- Collection: `GET /api/v1/invoices`, `POST /api/v1/invoices`
- Item: `GET /api/v1/invoices/{id}`, `PUT /api/v1/invoices/{id}`, `DELETE /api/v1/invoices/{id}`
- Special: `GET /api/v1/invoices/count` (query operations)
- Query params: Support filtering via `?client_id={id}`, `?status={status}`, `?q={search}`, `?sort={field}`, `?order={asc|desc}`
  - Implemented in `internal/query/query_params.go` and `internal/query/invoice_filters.go`
  - Example: `/api/v1/invoices?status=draft&sort=due_date&order=asc`

## Type System Notes

### Invoice Flow (Frontend → Backend)

1. **Create/Edit**: `invoice-editor.svelte` → form state → `Invoice` type
2. **Dates**: Issue date & due date as `DateValue` → synced to string via `$effect`
3. **Line Items**: Date picker per item, synced via `lineItemDateValues` array
4. **Save**: POST/PUT to `/api/v1/invoices` with complete `Invoice` object

### Party (Client/Provider) Flow

1. **Selection**: `profile-selector.svelte` → dropdown selection
2. **Display**: Read-only `Party` object with `name`, `email`, `abn`, etc.
3. **Storage**: Saved in `db/clients/` and `db/providers/`

## Important Conventions

- **No `@internationalized/date` imports in utility files** — use directly in `.svelte` files only
- **$effect for runes** — NEVER in `.ts` files, always in `.svelte`
- **Date strings**: Always ISO YYYY-MM-DD format (use `.toISOString().split('T')[0]`)
- **API base**: `/api/v1/` (no trailing slash on collection endpoints)
- **CORS**: Enabled for all origins (configured in `main.go`)
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

1. **Binary finds no data**: Check binary is run from directory containing `db/` folder
2. **Date picker closed immediately**: Ensure `Popover.Root bind:open` is implemented
3. **Props not reactive**: Always use `$props()` in component scripts, not destructuring imports
4. **Date sync loops**: Check that sync conditions prevent infinite updates (compare `.toString()`)
