/**
 * Store exports - Single point of import for all stores
 *
 * Usage:
 * import { invoices, loadInvoices, saveInvoice } from '@/stores';
 */

// Auth store
export { authStore, isAuthenticated, currentUserEmail, requiresOAuth } from './auth';

// Provider store
export { providers, activeProvider, providersLoading } from './provider';

// Client store
export {
	clients,
	clientsLoading,
	searchClients,
	addClient,
	updateClient,
	removeClient
} from './clients';

// Invoice store
export {
	invoices,
	invoicesLoading,
	invoicesError,
	invoiceFilters,
	filteredInvoices
} from './invoices';
