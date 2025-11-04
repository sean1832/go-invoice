/**
 * Store exports - Single point of import for all stores
 *
 * Usage:
 * import { invoices, loadInvoices, saveInvoice } from '@/stores';
 */

// Provider store
export {
	providers,
	activeProvider,
	providersLoading,
} from './provider';

// Client store
export {
	clients,
	clientsLoading,
	searchClients
} from './clients';

// Invoice store
export {
	invoices,
	invoicesLoading,
	invoicesError,
	invoiceFilters,
	filteredInvoices,
	generateInvoiceId
} from './invoices';

// Mock data (for development only - delete when switching to API)
export { mockProviders, mockClients, mockInvoices } from './mockup';
