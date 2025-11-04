import { get } from 'svelte/store';
import type { PageLoad } from '../$types';
import { invoices } from '@/stores';
import { api } from '@/services';
export const prerender = false;

export const load: PageLoad = async ({ params, fetch }) => {
	// try get from store first
	const allInvoices = get(invoices);
	let invoice = allInvoices.find((inv) => inv.id === params.id);

	if (!invoice && allInvoices.length === 0) {
		try {
			invoice = await api.invoices.getInvoice(fetch, params.id);
		} catch (error) {
			console.error(`failed to load invoice ${params.id}: `, error);
			return {
				invoice: null,
				error: error instanceof Error ? error.message : 'failed to load invoice data'
			};
		}
	}

	if (!invoice) {
		return {
			invoice: null,
			error: 'invoice not found'
		};
	}

	return {
		invoice
	};
};
