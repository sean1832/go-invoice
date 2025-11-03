/**
 * Client Store - Handles all client data operations
 * This is the single source of truth for client data in the frontend
 *
 * When switching to API:
 * 1. Replace localStorage operations with fetch calls to /api/v1/clients
 * 2. Keep the same function signatures
 * 3. Components won't need any changes
 */

import { writable } from 'svelte/store';
import type { ClientData } from '@/types/invoice';
import { mockClients } from './mockup';

// Client cache - this is what components will subscribe to
export const clients = writable<ClientData[]>([]);

// Loading state
export const clientsLoading = writable<boolean>(false);

// Initialize the store
let initialized = false;

/**
 * Load all clients from storage/API
 * Call this once when the app starts or when you need to refresh
 */
export async function loadClients(): Promise<ClientData[]> {
	if (typeof window === 'undefined') return [];

	clientsLoading.set(true);

	try {
		// TODO: Replace with API call
		// const response = await fetch('/api/v1/clients');
		// const data = await response.json();

		// Mock implementation using localStorage
		const stored = localStorage.getItem('clients');
		let clientList: ClientData[];

		if (stored) {
			clientList = JSON.parse(stored);
		} else {
			// Initialize with mock data on first load
			clientList = mockClients;
			localStorage.setItem('clients', JSON.stringify(clientList));
		}

		clients.set(clientList);
		initialized = true;
		return clientList;
	} catch (error) {
		console.error('Failed to load clients:', error);
		return [];
	} finally {
		clientsLoading.set(false);
	}
}

/**
 * Get a single client by ID
 */
export async function getClient(id: string): Promise<ClientData | null> {
	// Ensure clients are loaded
	if (!initialized) {
		await loadClients();
	}

	// TODO: Replace with API call
	// const response = await fetch(`/api/v1/clients/${id}`);
	// return await response.json();

	// Mock implementation
	const stored = localStorage.getItem('clients');
	if (!stored) return null;

	const clientList: ClientData[] = JSON.parse(stored);
	return clientList.find((c) => c.id === id) || null;
}

/**
 * Save a client (create or update)
 */
export async function saveClient(client: ClientData): Promise<void> {
	if (typeof window === 'undefined') return;

	clientsLoading.set(true);

	try {
		// TODO: Replace with API call
		// const method = client.id ? 'PUT' : 'POST';
		// const url = client.id ? `/api/v1/clients/${client.id}` : '/api/v1/clients';
		// await fetch(url, {
		// 	method,
		// 	headers: { 'Content-Type': 'application/json' },
		// 	body: JSON.stringify(client)
		// });

		// Mock implementation
		const stored = localStorage.getItem('clients');
		const clientList: ClientData[] = stored ? JSON.parse(stored) : [];

		const index = clientList.findIndex((c) => c.id === client.id);
		if (index >= 0) {
			clientList[index] = client;
		} else {
			clientList.push(client);
		}

		localStorage.setItem('clients', JSON.stringify(clientList));

		// Update the store
		clients.set(clientList);
	} catch (error) {
		console.error('Failed to save client:', error);
		throw error;
	} finally {
		clientsLoading.set(false);
	}
}

/**
 * Delete a client
 */
export async function deleteClient(id: string): Promise<void> {
	if (typeof window === 'undefined') return;

	clientsLoading.set(true);

	try {
		// TODO: Replace with API call
		// await fetch(`/api/v1/clients/${id}`, { method: 'DELETE' });

		// Mock implementation
		const stored = localStorage.getItem('clients');
		if (!stored) return;

		const clientList: ClientData[] = JSON.parse(stored);
		const filtered = clientList.filter((c) => c.id !== id);

		localStorage.setItem('clients', JSON.stringify(filtered));

		// Update the store
		clients.set(filtered);
	} catch (error) {
		console.error('Failed to delete client:', error);
		throw error;
	} finally {
		clientsLoading.set(false);
	}
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
