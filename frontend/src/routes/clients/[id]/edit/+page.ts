import { get } from 'svelte/store';
import type { PageLoad } from './$types';
import { clients } from '@/stores';
import { api } from '@/services';

export const prerender = false;

export const load: PageLoad = async ({ params, fetch }) => {
	// try get from store first

	const allClients = get(clients);
	let client = allClients.find((cli) => cli.id === params.id);

	if (!client) {
		try {
			client = await api.clients.getClient(fetch, params.id);
		} catch (error) {
			console.error(`failed to load client ${params.id}: `, error);
			return {
				client: null,
				error: error instanceof Error ? error.message : 'failed to load client data'
			};
		}
	}

	if (!client) {
		return {
			client: null,
			error: 'client not found'
		};
	}

	return {
		client
	};
};
