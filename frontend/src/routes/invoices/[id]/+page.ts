import { invoices } from '$lib/stores/invoices';
import { api } from '$lib/services';
import { get } from 'svelte/store';
import type { PageLoad } from './$types';

export const prerender = false;

/**
 * Load function to get invoice from store or fetch from API
 * Tries store first (if populated), falls back to API call
 */
export const load: PageLoad = async ({ params, fetch }) => {
	// Try to get from store first
	const allInvoices = get(invoices);
	let invoice = allInvoices.find((inv) => inv.id === params.id);

	// If not in store, fetch from API
	if (!invoice && allInvoices.length === 0) {
		try {
			invoice = await api.invoices.getInvoice(fetch, params.id);
		} catch (error) {
			console.error(`Failed to load invoice ${params.id}:`, error);
			return {
				invoice: null,
				error: error instanceof Error ? error.message : 'Failed to load invoice'
			};
		}
	}

	if (!invoice) {
		return {
			invoice: null,
			error: 'Invoice not found'
		};
	}

	return {
		invoice
	};
};
