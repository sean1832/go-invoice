/**
 * Provider Store - Handles all provider data operations
 * This is the single source of truth for provider data in the frontend
 *
 * When switching to API:
 * 1. Replace localStorage operations with fetch calls to /api/v1/providers
 * 2. Keep the same function signatures
 * 3. Components won't need any changes
 */

import { writable } from 'svelte/store';
import type { ProviderData } from '@/types/invoice';
import { mockProviders } from './mockup';

// Provider cache - this is what components will subscribe to
export const providers = writable<ProviderData[]>([]);

// Active provider - the currently selected provider for creating invoices
export const activeProvider = writable<ProviderData | null>(null);

// Loading state
export const providersLoading = writable<boolean>(false);

// Initialize stores with localStorage data
let initialized = false;

function initializeActiveProvider(): void {
	if (typeof window === 'undefined') return;

	const stored = localStorage.getItem('activeProvider');
	if (stored) {
		try {
			activeProvider.set(JSON.parse(stored));
		} catch (error) {
			console.error('Failed to parse active provider from localStorage:', error);
		}
	}
}

// Subscribe to save active provider changes to localStorage
if (typeof window !== 'undefined') {
	activeProvider.subscribe((value) => {
		if (value) {
			localStorage.setItem('activeProvider', JSON.stringify(value));
		} else {
			localStorage.removeItem('activeProvider');
		}
	});
}

/**
 * Load all providers from storage/API
 * Call this once when the app starts or when you need to refresh
 */
export async function loadProviders(): Promise<ProviderData[]> {
	if (typeof window === 'undefined') return [];

	providersLoading.set(true);

	try {
		// TODO: Replace with API call
		// const response = await fetch('/api/v1/providers');
		// const data = await response.json();

		// Mock implementation using localStorage
		const stored = localStorage.getItem('providers');
		let providerList: ProviderData[];

		if (stored) {
			providerList = JSON.parse(stored);
		} else {
			// Initialize with mock data on first load
			providerList = mockProviders;
			localStorage.setItem('providers', JSON.stringify(providerList));
		}

		providers.set(providerList);

		// If no active provider is set, set the first one
		if (!initialized) {
			initializeActiveProvider();
			const current = await new Promise<ProviderData | null>((resolve) => {
				const unsubscribe = activeProvider.subscribe((value) => {
					resolve(value);
					unsubscribe();
				});
			});

			if (!current && providerList.length > 0) {
				setActiveProvider(providerList[0]);
			}
			initialized = true;
		}

		return providerList;
	} catch (error) {
		console.error('Failed to load providers:', error);
		return [];
	} finally {
		providersLoading.set(false);
	}
}

/**
 * Get a single provider by ID
 */
export async function getProvider(id: string): Promise<ProviderData | null> {
	// Ensure providers are loaded
	if (!initialized) {
		await loadProviders();
	}

	// TODO: Replace with API call
	// const response = await fetch(`/api/v1/providers/${id}`);
	// return await response.json();

	// Mock implementation
	const stored = localStorage.getItem('providers');
	if (!stored) return null;

	const providerList: ProviderData[] = JSON.parse(stored);
	return providerList.find((p) => p.id === id) || null;
}

/**
 * Save a provider (create or update)
 */
export async function saveProvider(provider: ProviderData): Promise<void> {
	if (typeof window === 'undefined') return;

	providersLoading.set(true);

	try {
		// TODO: Replace with API call
		// const method = provider.id ? 'PUT' : 'POST';
		// const url = provider.id ? `/api/v1/providers/${provider.id}` : '/api/v1/providers';
		// await fetch(url, {
		// 	method,
		// 	headers: { 'Content-Type': 'application/json' },
		// 	body: JSON.stringify(provider)
		// });

		// Mock implementation
		const stored = localStorage.getItem('providers');
		const providerList: ProviderData[] = stored ? JSON.parse(stored) : [];

		const index = providerList.findIndex((p) => p.id === provider.id);
		if (index >= 0) {
			providerList[index] = provider;
		} else {
			providerList.push(provider);
		}

		localStorage.setItem('providers', JSON.stringify(providerList));

		// Update the store
		providers.set(providerList);

		// If this is the active provider or no active provider exists, update it
		const current = await new Promise<ProviderData | null>((resolve) => {
			const unsubscribe = activeProvider.subscribe((value) => {
				resolve(value);
				unsubscribe();
			});
		});

		if (!current || current.id === provider.id) {
			setActiveProvider(provider);
		}
	} catch (error) {
		console.error('Failed to save provider:', error);
		throw error;
	} finally {
		providersLoading.set(false);
	}
}

/**
 * Delete a provider
 */
export async function deleteProvider(id: string): Promise<void> {
	if (typeof window === 'undefined') return;

	providersLoading.set(true);

	try {
		// TODO: Replace with API call
		// await fetch(`/api/v1/providers/${id}`, { method: 'DELETE' });

		// Mock implementation
		const stored = localStorage.getItem('providers');
		if (!stored) return;

		const providerList: ProviderData[] = JSON.parse(stored);
		const filtered = providerList.filter((p) => p.id !== id);

		localStorage.setItem('providers', JSON.stringify(filtered));

		// Update the store
		providers.set(filtered);

		// If the deleted provider was active, clear or set a new one
		const current = await new Promise<ProviderData | null>((resolve) => {
			const unsubscribe = activeProvider.subscribe((value) => {
				resolve(value);
				unsubscribe();
			});
		});

		if (current?.id === id) {
			if (filtered.length > 0) {
				setActiveProvider(filtered[0]);
			} else {
				clearActiveProvider();
			}
		}
	} catch (error) {
		console.error('Failed to delete provider:', error);
		throw error;
	} finally {
		providersLoading.set(false);
	}
}

/**
 * Set the active provider (used for creating new invoices)
 */
export function setActiveProvider(provider: ProviderData): void {
	activeProvider.set(provider);
}

/**
 * Clear the active provider
 */
export function clearActiveProvider(): void {
	activeProvider.set(null);
}

/**
 * Initialize providers on app startup
 * Call this from your root layout
 */
export async function initializeProviders(): Promise<void> {
	await loadProviders();
}
