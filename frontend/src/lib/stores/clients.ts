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

export function addClient(item: ClientData): void {
	clients.update((items) => {
		// Ensure items is always an array
		const itemsArray = Array.isArray(items) ? items : [];

		// Check if client already exists
		const exists = itemsArray.some((c) => c.id === item.id);
		if (exists) {
			// Update existing
			return itemsArray.map((c) => (c.id === item.id ? { ...c, ...item } : c));
		} else {
			// Add new
			return [...itemsArray, item];
		}
	});
}

export function updateClient(item: ClientData): void {
	clients.update((items) => {
		const itemsArray = Array.isArray(items) ? items : [];
		return itemsArray.map((c) => (c.id === item.id ? { ...c, ...item } : c));
	});
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
