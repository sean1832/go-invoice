# Store Architecture Guide

## Overview

The frontend uses a **centralized store architecture** where all data operations (CRUD) are handled through Svelte stores. This approach provides:

- ✅ **Single source of truth** for data
- ✅ **Easy API integration** - just swap localStorage with fetch calls
- ✅ **Reactive UI** - components auto-update when data changes
- ✅ **Clean separation** - components don't know about data fetching
- ✅ **Easy testing** - stores can be mocked independently

## Store Structure

```
frontend/src/lib/stores/
├── index.ts          # Barrel export - import everything from here
├── mockup.ts         # Mock data (DELETE when switching to API)
├── provider.ts       # Provider CRUD + active provider
├── clients.ts        # Client CRUD
└── invoices.ts       # Invoice CRUD + filtering
```

## Usage Pattern

### 1. Import from stores

```typescript
import {
  invoices,
  loadInvoices,
  saveInvoice,
  deleteInvoice
} from '@/stores';
```

### 2. Load data (usually in root layout or page load)

```typescript
import { onMount } from 'svelte';

onMount(async () => {
  await loadInvoices();
});
```

### 3. Subscribe to store in components

```svelte
<script lang="ts">
	import { invoices } from '@/stores';
</script>

{#each $invoices as invoice}
	<InvoiceCard {invoice} />
{/each}
```

### 4. Modify data through store functions

```typescript
// Create/Update
await saveInvoice(newInvoice);

// Delete
await deleteInvoice('INV-001');

// The UI automatically updates!
```

## Store API Reference

### Provider Store

```typescript
// Stores
providers              // writable<ProviderData[]> - all providers
activeProvider         // writable<ProviderData | null> - current provider
providersLoading       // writable<boolean> - loading state

// Functions
loadProviders()        // Load all providers
getProvider(id)        // Get single provider
saveProvider(provider) // Create/update provider
deleteProvider(id)     // Delete provider
setActiveProvider(p)   // Set active provider
clearActiveProvider()  // Clear active provider
initializeProviders()  // Initialize on app start
```

### Client Store

```typescript
// Stores
clients                // writable<ClientData[]> - all clients
clientsLoading         // writable<boolean> - loading state

// Functions
loadClients()          // Load all clients
getClient(id)          // Get single client
saveClient(client)     // Create/update client
deleteClient(id)       // Delete client
searchClients(query)   // Search clients by name/email/abn
```

### Invoice Store

```typescript
// Stores
invoices               // writable<Invoice[]> - all invoices
invoicesLoading        // writable<boolean> - loading state
invoiceFilters         // writable<InvoiceFilters> - current filters
filteredInvoices       // derived store - filtered & sorted invoices

// Functions
loadInvoices()         // Load all invoices
getInvoice(id)         // Get single invoice
saveInvoice(invoice)   // Create/update invoice
deleteInvoice(id)      // Delete invoice
updateInvoiceStatus(id, status) // Update status only
getInvoiceCount(status?) // Get count by status
generateInvoiceId()    // Generate new ID
```

## Switching from Mock to API

### Current Implementation (Mock)

```typescript
export async function loadInvoices(): Promise<Invoice[]> {
  // Mock implementation using localStorage
  const stored = localStorage.getItem('invoices');
  let invoiceList: Invoice[];

  if (stored) {
    invoiceList = JSON.parse(stored);
  } else {
    invoiceList = mockInvoices;
    localStorage.setItem('invoices', JSON.stringify(invoiceList));
  }

  invoices.set(invoiceList);
  return invoiceList;
}
```

### API Implementation (Future)

```typescript
export async function loadInvoices(): Promise<Invoice[]> {
  invoicesLoading.set(true);

  try {
    const response = await fetch('/api/v1/invoices');
    if (!response.ok) throw new Error('Failed to load invoices');

    const data = await response.json();
    invoices.set(data);
    return data;
  } catch (error) {
    console.error('Failed to load invoices:', error);
    return [];
  } finally {
    invoicesLoading.set(false);
  }
}
```

**That's it!** No changes needed in components.

## API Migration Checklist

When you're ready to switch to the real API:

1. ✅ **Ensure backend API is running** (`http://localhost:8080/api/v1/`)

2. ✅ **Update each store file** (provider.ts, clients.ts, invoices.ts):
   - Replace localStorage logic with fetch calls
   - Add proper error handling
   - Keep the same function signatures

3. ✅ **Delete mock data**:

   ```bash
   rm frontend/src/lib/stores/mockup.ts
   ```

   - Remove `mockup.ts` export from `index.ts`

4. ✅ **Test each operation**:
   - Load data
   - Create new items
   - Update existing items
   - Delete items
   - Filter/search

5. ✅ **Update error handling** in components if needed

## Best Practices

### ✅ DO

- **Load data once** in root layout or page component
- **Subscribe with $store** syntax in components
- **Call store functions** for all data operations
- **Use derived stores** for computed/filtered data
- **Handle loading states** with `invoicesLoading`, etc.

### ❌ DON'T

- **Don't fetch data in components** - use store functions
- **Don't duplicate data** - stores are the single source of truth
- **Don't mutate store data directly** - use store functions
- **Don't mix mock and API** - choose one approach
- **Don't forget to await** store functions (they're async)

## Example: Complete Invoice Flow

```svelte
<!-- routes/invoices/new/+page.svelte -->
<script lang="ts">
	import { saveInvoice, generateInvoiceId } from '@/stores';
	import InvoiceEditor from '@/components/custom/invoice-editor/invoice-editor.svelte';

	async function handleSave(data: any) {
		try {
			await saveInvoice(data);
			// Redirect to invoice list
			window.location.href = '/';
		} catch (error) {
			console.error('Failed to save invoice:', error);
			alert('Failed to save invoice. Please try again.');
		}
	}
</script>

<InvoiceEditor mode="create" onSave={handleSave} />
```

## Derived Stores Example

The invoice store includes a `filteredInvoices` derived store that automatically updates when filters change:

```typescript
export const filteredInvoices = derived(
  [invoices, invoiceFilters],
  ([$invoices, $filters]) => {
    let result = [...$invoices];

    // Apply filters
    if ($filters.status && $filters.status !== 'all') {
      result = result.filter((inv) => inv.status === $filters.status);
    }

    if ($filters.searchQuery) {
      const query = $filters.searchQuery.toLowerCase();
      result = result.filter((inv) =>
        inv.id.toLowerCase().includes(query) ||
        inv.client.name.toLowerCase().includes(query)
      );
    }

    // Apply sorting
    result.sort((a, b) => {
      // ... sorting logic
    });

    return result;
  }
);
```

Usage:

```svelte
<script>
	import { filteredInvoices, invoiceFilters } from '@/stores';
</script>

<!-- Update filter -->
<button onclick={() => invoiceFilters.set({ status: 'draft' })}> Show Drafts </button>

<!-- Display filtered results -->
{#each $filteredInvoices as invoice}
	<InvoiceCard {invoice} />
{/each}
```

## Loading States

All stores expose loading states for UI feedback:

```svelte
<script>
	import { invoices, invoicesLoading, loadInvoices } from '@/stores';
	import { onMount } from 'svelte';

	onMount(() => loadInvoices());
</script>

{#if $invoicesLoading}
	<LoadingSpinner />
{:else}
	{#each $invoices as invoice}
		<InvoiceCard {invoice} />
	{/each}
{/if}
```

## Troubleshooting

### Store data is empty

- Make sure you call `loadInvoices()`, `loadClients()`, etc. in your root layout
- Check browser console for errors
- Verify localStorage has data (DevTools → Application → Local Storage)

### Updates not reflecting in UI

- Ensure you're subscribing with `$store` syntax
- Verify you're calling store functions (not mutating directly)
- Check that store functions call `.set()` or `.update()`

### API calls not working

- Verify backend is running on `http://localhost:8080`
- Check network tab in DevTools
- Ensure CORS is properly configured in backend
- Verify request/response formats match backend expectations

## Related Files

- **Types**: `frontend/src/lib/types/invoice.ts`
- **Components**: `frontend/src/lib/components/custom/`
- **Root Layout**: `frontend/src/routes/+layout.svelte` (initializes stores)
