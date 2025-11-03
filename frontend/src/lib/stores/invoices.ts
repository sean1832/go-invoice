/**
 * Invoice Store - Handles all invoice data operations
 * This is the single source of truth for invoice data in the frontend
 *
 * When switching to API:
 * 1. Replace localStorage operations with fetch calls to /api/v1/invoices
 * 2. Keep the same function signatures
 * 3. Components won't need any changes
 */

import { writable, derived } from 'svelte/store';
import type { Invoice, InvoiceFilters, InvoiceStatus } from '@/types/invoice';
import { mockInvoices } from './mockup';

// Invoice cache - this is what components will subscribe to
export const invoices = writable<Invoice[]>([]);

// Loading state
export const invoicesLoading = writable<boolean>(false);

// Current filters
export const invoiceFilters = writable<InvoiceFilters>({
	status: 'all',
	searchQuery: '',
	sortBy: 'date',
	sortOrder: 'desc'
});

// Initialize invoices synchronously from localStorage
if (typeof window !== 'undefined') {
	const stored = localStorage.getItem('invoices');
	if (stored) {
		try {
			invoices.set(JSON.parse(stored));
		} catch (error) {
			console.error('Failed to parse invoices from localStorage:', error);
		}
	} else {
		// Initialize with mock data if nothing in storage
		invoices.set(mockInvoices);
		localStorage.setItem('invoices', JSON.stringify(mockInvoices));
	}
}

// Derived store for filtered invoices
export const filteredInvoices = derived([invoices, invoiceFilters], ([$invoices, $filters]) => {
	let result = [...$invoices];

	// Filter by status
	if ($filters.status && $filters.status !== 'all') {
		result = result.filter((inv) => inv.status === $filters.status);
	}

	// Filter by search query
	if ($filters.searchQuery) {
		const query = $filters.searchQuery.toLowerCase();
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

		switch ($filters.sortBy) {
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

		return $filters.sortOrder === 'asc' ? comparison : -comparison;
	});

	return result;
});

// Initialize the store
let initialized = false;

/**
 * Load all invoices from storage/API
 * Call this once when the app starts or when you need to refresh
 */
export async function loadInvoices(): Promise<Invoice[]> {
	if (typeof window === 'undefined') return [];

	invoicesLoading.set(true);

	try {
		// TODO: Replace with API call
		// const response = await fetch('/api/v1/invoices');
		// const data = await response.json();

		// Mock implementation using localStorage
		const stored = localStorage.getItem('invoices');
		let invoiceList: Invoice[];

		if (stored) {
			invoiceList = JSON.parse(stored);
		} else {
			// Initialize with mock data on first load
			invoiceList = mockInvoices;
			localStorage.setItem('invoices', JSON.stringify(invoiceList));
		}

		invoices.set(invoiceList);
		initialized = true;
		return invoiceList;
	} catch (error) {
		console.error('Failed to load invoices:', error);
		return [];
	} finally {
		invoicesLoading.set(false);
	}
}

/**
 * Get a single invoice by ID
 */
export async function getInvoice(id: string): Promise<Invoice | null> {
	// Ensure invoices are loaded
	if (!initialized) {
		await loadInvoices();
	}

	// TODO: Replace with API call
	// const response = await fetch(`/api/v1/invoices/${id}`);
	// return await response.json();

	// Mock implementation
	const stored = localStorage.getItem('invoices');
	if (!stored) return null;

	const invoiceList: Invoice[] = JSON.parse(stored);
	return invoiceList.find((inv) => inv.id === id) || null;
}

/**
 * Save an invoice (create or update)
 */
export async function saveInvoice(invoice: Invoice): Promise<void> {
	if (typeof window === 'undefined') return;

	invoicesLoading.set(true);

	try {
		// TODO: Replace with API call
		// const method = invoice.id ? 'PUT' : 'POST';
		// const url = invoice.id ? `/api/v1/invoices/${invoice.id}` : '/api/v1/invoices';
		// await fetch(url, {
		// 	method,
		// 	headers: { 'Content-Type': 'application/json' },
		// 	body: JSON.stringify(invoice)
		// });

		// Mock implementation
		const stored = localStorage.getItem('invoices');
		const invoiceList: Invoice[] = stored ? JSON.parse(stored) : [];

		const index = invoiceList.findIndex((inv) => inv.id === invoice.id);
		if (index >= 0) {
			invoiceList[index] = invoice;
		} else {
			invoiceList.push(invoice);
		}

		localStorage.setItem('invoices', JSON.stringify(invoiceList));

		// Update the store
		invoices.set(invoiceList);
	} catch (error) {
		console.error('Failed to save invoice:', error);
		throw error;
	} finally {
		invoicesLoading.set(false);
	}
}

/**
 * Delete an invoice
 */
export async function deleteInvoice(id: string): Promise<void> {
	if (typeof window === 'undefined') return;

	invoicesLoading.set(true);

	try {
		// TODO: Replace with API call
		// await fetch(`/api/v1/invoices/${id}`, { method: 'DELETE' });

		// Mock implementation
		const stored = localStorage.getItem('invoices');
		if (!stored) return;

		const invoiceList: Invoice[] = JSON.parse(stored);
		const filtered = invoiceList.filter((inv) => inv.id !== id);

		localStorage.setItem('invoices', JSON.stringify(filtered));

		// Update the store
		invoices.set(filtered);
	} catch (error) {
		console.error('Failed to delete invoice:', error);
		throw error;
	} finally {
		invoicesLoading.set(false);
	}
}

/**
 * Update invoice status
 */
export async function updateInvoiceStatus(id: string, status: InvoiceStatus): Promise<void> {
	const invoice = await getInvoice(id);
	if (!invoice) {
		throw new Error(`Invoice ${id} not found`);
	}

	invoice.status = status;
	await saveInvoice(invoice);
}

/**
 * Get invoice count by status
 */
export async function getInvoiceCount(status?: InvoiceStatus): Promise<number> {
	// TODO: Replace with API call
	// const url = status ? `/api/v1/invoices/count?status=${status}` : '/api/v1/invoices/count';
	// const response = await fetch(url);
	// return await response.json();

	// Mock implementation
	const stored = localStorage.getItem('invoices');
	if (!stored) return 0;

	const invoiceList: Invoice[] = JSON.parse(stored);

	if (!status) return invoiceList.length;

	return invoiceList.filter((inv) => inv.status === status).length;
}

/**
 * Generate a new invoice ID based on current date
 */
export function generateInvoiceId(): string {
	const now = new Date();
	const year = now.getFullYear().toString().slice(-2);
	const month = (now.getMonth() + 1).toString().padStart(2, '0');
	const day = now.getDate().toString().padStart(2, '0');

	// Get count of invoices created today to generate a sequence number
	const prefix = `INV-${year}${month}${day}`;

	// TODO: This should ideally come from the backend to avoid conflicts
	const stored = localStorage.getItem('invoices');
	const invoiceList: Invoice[] = stored ? JSON.parse(stored) : [];
	const todayInvoices = invoiceList.filter((inv) => inv.id.startsWith(prefix));
	const sequence = (todayInvoices.length + 1).toString().padStart(3, '0');

	return `${prefix}${sequence}`;
}
