/**
 * Invoice Store - Shared state for invoice data across components
 */

import { writable, derived } from 'svelte/store';
import type { Invoice, InvoiceFilters } from '@/types/invoice';

// Stores
export const invoices = writable<Invoice[]>([]);
export const invoicesLoading = writable<boolean>(false);
export const invoicesError = writable<string | null>(null);

export const invoiceFilters = writable<InvoiceFilters>({
	status: 'all',
	search_query: '',
	sort_by: 'date',
	sort_order: 'desc'
});

// Filtered and sorted invoices based on current filters
export const filteredInvoices = derived([invoices, invoiceFilters], ([$invoices, $filters]) => {
	let result = [...$invoices];

	// Filter by status
	if ($filters.status && $filters.status !== 'all') {
		result = result.filter((inv) => inv.status === $filters.status);
	}

	// Filter by search query
	if ($filters.search_query) {
		const query = $filters.search_query.toLowerCase();
		result = result.filter(
			(inv) =>
				inv.id.toLowerCase().includes(query) ||
				inv.client.name.toLowerCase().includes(query) ||
				inv.client.email?.toLowerCase().includes(query)
		);
	}

	// Sort
	result.sort((a, b) => {
		let comparison = 0;

		switch ($filters.sort_by) {
			case 'date':
				comparison = new Date(a.date).getTime() - new Date(b.date).getTime();
				break;
			case 'amount':
				comparison = a.pricing.total - b.pricing.total;
				break;
			case 'client':
				comparison = a.client.name.localeCompare(b.client.name);
				break;
		}

		return $filters.sort_order === 'asc' ? comparison : -comparison;
	});

	return result;
});

/**
 * Remove an invoice from the store by ID (memory only, ensure to sync with backend)
 */
export function removeInvoice(id: string): void {
	invoices.update((items) => items.filter((inv) => inv.id !== id));
}