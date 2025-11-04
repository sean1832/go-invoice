/**
 * Client Store - Shared state for client data across components
 * This is the single source of truth for client data in the frontend
 */

import { writable } from 'svelte/store';
import type { ClientData } from '@/types/invoice';
// import { mockClients } from './mockup';

// Client cache - this is what components will subscribe to
export const clients = writable<ClientData[]>([]);

// Loading state
export const clientsLoading = writable<boolean>(false);

// Initialize the store
export const initialized = writable<boolean>(false);

export function removeClient(id: string): void {
	clients.update((items) => items.filter((inv) => inv.id !== id));
}

/**
 * Search clients by query
 */
export function searchClients(query: string): ClientData[] {
	let clientList: ClientData[] = [];
	clients.subscribe((value) => {
		clientList = value;
	})();

	if (!query) return clientList;

	const lowerQuery = query.toLowerCase();
	return clientList.filter(
		(client) =>
			client.name.toLowerCase().includes(lowerQuery) ||
			client.email?.toLowerCase().includes(lowerQuery) ||
			client.abn?.includes(lowerQuery)
	);
}
