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
 * Generate a new invoice ID based on current date
 */
export function generateInvoiceId(): string {
	const now = new Date();
	const year = now.getFullYear().toString().slice(-2);
	const month = (now.getMonth() + 1).toString().padStart(2, '0');
	const day = now.getDate().toString().padStart(2, '0');
	const prefix = `INV-${year}${month}${day}`;

	// TODO: This should come from the backend to avoid conflicts
	const stored = localStorage.getItem('invoices');
	const invoiceList: Invoice[] = stored ? JSON.parse(stored) : [];
	const todayInvoices = invoiceList.filter((inv) => inv.id.startsWith(prefix));
	const sequence = (todayInvoices.length + 1).toString().padStart(3, '0');

	return `${prefix}${sequence}`;
}
