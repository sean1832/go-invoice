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
// import { mockProviders } from './mockup';

// Provider cache - this is what components will subscribe to
export const providers = writable<ProviderData[]>([]);

// Active provider - the currently selected provider for creating invoices
export const activeProvider = writable<ProviderData | null>(null);

// Loading state
export const providersLoading = writable<boolean>(false);

// Initialize stores with localStorage data
export const initialized = writable<boolean>(false);

// Initialize activeProvider from localStorage on startup
if (typeof window !== 'undefined') {
	const storedActiveProvider = localStorage.getItem('activeProvider');
	if (storedActiveProvider) {
		try {
			activeProvider.set(JSON.parse(storedActiveProvider));
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

// // Initialize providers synchronously from localStorage
// if (typeof window !== 'undefined') {
// 	const stored = localStorage.getItem('providers');
// 	if (stored) {
// 		try {
// 			const storedProviders = JSON.parse(stored);
// 			// Check if localStorage data matches mock data count
// 			// This helps during development when mock data changes
// 			if (storedProviders.length !== mockProviders.length) {
// 				console.warn(
// 					`[Provider Store] localStorage has ${storedProviders.length} providers, but mock has ${mockProviders.length}. Syncing with mock data...`
// 				);
// 				// Merge: keep mock data and add any custom providers from localStorage
// 				const mockIds = new Set(mockProviders.map((p) => p.id));
// 				const customProviders = storedProviders.filter((p: ProviderData) => !mockIds.has(p.id));
// 				const mergedProviders = [...mockProviders, ...customProviders];
// 				providers.set(mergedProviders);
// 				localStorage.setItem('providers', JSON.stringify(mergedProviders));
// 			} else {
// 				providers.set(storedProviders);
// 			}
// 		} catch (error) {
// 			console.error('Failed to parse providers from localStorage:', error);
// 			// Fallback to mock data on error
// 			providers.set(mockProviders);
// 			localStorage.setItem('providers', JSON.stringify(mockProviders));
// 		}
// 	} else {
// 		// Initialize with mock data if nothing in storage
// 		providers.set(mockProviders);
// 		localStorage.setItem('providers', JSON.stringify(mockProviders));
// 	}
// }

// function initializeActiveProvider(): void {
// 	if (typeof window === 'undefined') return;

// 	const stored = localStorage.getItem('activeProvider');
// 	if (stored) {
// 		try {
// 			activeProvider.set(JSON.parse(stored));
// 		} catch (error) {
// 			console.error('Failed to parse active provider from localStorage:', error);
// 		}
// 	}
// }

// // Initialize activeProvider immediately from localStorage (synchronous)
// if (typeof window !== 'undefined') {
// 	initializeActiveProvider();
// }

// // Subscribe to save active provider changes to localStorage
// if (typeof window !== 'undefined') {
// 	activeProvider.subscribe((value) => {
// 		console.log(
// 			'[Provider Store] activeProvider subscriber triggered with:',
// 			value?.name || 'null'
// 		);
// 		if (value) {
// 			localStorage.setItem('activeProvider', JSON.stringify(value));
// 			console.log('[Provider Store] Saved to localStorage:', value.name);
// 		} else {
// 			localStorage.removeItem('activeProvider');
// 			console.log('[Provider Store] Removed from localStorage');
// 		}
// 	});
// }

// /**
//  * Load all providers from storage/API
//  * Call this once when the app starts or when you need to refresh
//  */
// export async function loadProviders(): Promise<ProviderData[]> {
// 	if (typeof window === 'undefined') return [];

// 	providersLoading.set(true);
// 	console.log('[Provider Store] Loading providers...');

// 	try {
// 		// TODO: Replace with API call
// 		// const response = await fetch('/api/v1/providers');
// 		// const data = await response.json();

// 		// Mock implementation using localStorage
// 		const stored = localStorage.getItem('providers');
// 		let providerList: ProviderData[];

// 		if (stored) {
// 			const storedProviders = JSON.parse(stored);
// 			// Check if localStorage data matches mock data count
// 			if (storedProviders.length !== mockProviders.length) {
// 				console.log(
// 					`[Provider Store] localStorage has ${storedProviders.length} providers, mock has ${mockProviders.length}. Syncing...`
// 				);
// 				// Merge: keep mock data and add any custom providers from localStorage
// 				const mockIds = new Set(mockProviders.map((p) => p.id));
// 				const customProviders = storedProviders.filter((p: ProviderData) => !mockIds.has(p.id));
// 				providerList = [...mockProviders, ...customProviders];
// 				localStorage.setItem('providers', JSON.stringify(providerList));
// 				console.log('[Provider Store] Synced providers:', providerList.length, 'providers');
// 			} else {
// 				providerList = storedProviders;
// 				console.log('[Provider Store] Loaded from localStorage:', providerList.length, 'providers');
// 			}
// 		} else {
// 			// Initialize with mock data on first load
// 			providerList = mockProviders;
// 			localStorage.setItem('providers', JSON.stringify(providerList));
// 			console.log('[Provider Store] Initialized with mock data:', providerList.length, 'providers');
// 		}

// 		providers.set(providerList);

// 		// Validate and initialize active provider
// 		const current = await new Promise<ProviderData | null>((resolve) => {
// 			const unsubscribe = activeProvider.subscribe((value) => {
// 				resolve(value);
// 				unsubscribe();
// 			});
// 		});

// 		console.log('[Provider Store] Current active provider:', current?.name || 'none');

// 		// Check if current active provider exists in the loaded list
// 		const activeExists = current && providerList.some((p) => p.id === current.id);

// 		if (!activeExists && providerList.length > 0) {
// 			// If active provider doesn't exist or is null, set the first one
// 			console.log(
// 				'[Provider Store] Setting active provider to first in list:',
// 				providerList[0].name
// 			);
// 			setActiveProvider(providerList[0]);
// 		} else if (activeExists && current) {
// 			// Update the active provider with the latest data from the list
// 			const updatedProvider = providerList.find((p) => p.id === current.id);
// 			if (updatedProvider) {
// 				console.log(
// 					'[Provider Store] Updating active provider with latest data:',
// 					updatedProvider.name
// 				);
// 				setActiveProvider(updatedProvider);
// 			}
// 		}

// 		initialized = true;

// 		return providerList;
// 	} catch (error) {
// 		console.error('Failed to load providers:', error);
// 		return [];
// 	} finally {
// 		providersLoading.set(false);
// 	}
// }

// /**
//  * Get a single provider by ID
//  */
// export async function getProvider(id: string): Promise<ProviderData | null> {
// 	// Ensure providers are loaded
// 	if (!initialized) {
// 		await loadProviders();
// 	}

// 	// TODO: Replace with API call
// 	// const response = await fetch(`/api/v1/providers/${id}`);
// 	// return await response.json();

// 	// Mock implementation
// 	const stored = localStorage.getItem('providers');
// 	if (!stored) return null;

// 	const providerList: ProviderData[] = JSON.parse(stored);
// 	return providerList.find((p) => p.id === id) || null;
// }

// /**
//  * Save a provider (create or update)
//  */
// export async function saveProvider(provider: ProviderData): Promise<void> {
// 	if (typeof window === 'undefined') return;

// 	providersLoading.set(true);

// 	try {
// 		// TODO: Replace with API call
// 		// const method = provider.id ? 'PUT' : 'POST';
// 		// const url = provider.id ? `/api/v1/providers/${provider.id}` : '/api/v1/providers';
// 		// await fetch(url, {
// 		// 	method,
// 		// 	headers: { 'Content-Type': 'application/json' },
// 		// 	body: JSON.stringify(provider)
// 		// });

// 		// Mock implementation
// 		const stored = localStorage.getItem('providers');
// 		const providerList: ProviderData[] = stored ? JSON.parse(stored) : [];

// 		const index = providerList.findIndex((p) => p.id === provider.id);
// 		if (index >= 0) {
// 			providerList[index] = provider;
// 		} else {
// 			providerList.push(provider);
// 		}

// 		localStorage.setItem('providers', JSON.stringify(providerList));

// 		// Update the store
// 		providers.set(providerList);

// 		// If this is the active provider or no active provider exists, update it
// 		const current = await new Promise<ProviderData | null>((resolve) => {
// 			const unsubscribe = activeProvider.subscribe((value) => {
// 				resolve(value);
// 				unsubscribe();
// 			});
// 		});

// 		if (!current || current.id === provider.id) {
// 			setActiveProvider(provider);
// 		}
// 	} catch (error) {
// 		console.error('Failed to save provider:', error);
// 		throw error;
// 	} finally {
// 		providersLoading.set(false);
// 	}
// }

// /**
//  * Delete a provider
//  */
// export async function deleteProvider(id: string): Promise<void> {
// 	if (typeof window === 'undefined') return;

// 	providersLoading.set(true);

// 	try {
// 		// TODO: Replace with API call
// 		// await fetch(`/api/v1/providers/${id}`, { method: 'DELETE' });

// 		// Mock implementation
// 		const stored = localStorage.getItem('providers');
// 		if (!stored) return;

// 		const providerList: ProviderData[] = JSON.parse(stored);
// 		const filtered = providerList.filter((p) => p.id !== id);

// 		localStorage.setItem('providers', JSON.stringify(filtered));

// 		// Update the store
// 		providers.set(filtered);

// 		// If the deleted provider was active, clear or set a new one
// 		const current = await new Promise<ProviderData | null>((resolve) => {
// 			const unsubscribe = activeProvider.subscribe((value) => {
// 				resolve(value);
// 				unsubscribe();
// 			});
// 		});

// 		if (current?.id === id) {
// 			if (filtered.length > 0) {
// 				setActiveProvider(filtered[0]);
// 			} else {
// 				clearActiveProvider();
// 			}
// 		}
// 	} catch (error) {
// 		console.error('Failed to delete provider:', error);
// 		throw error;
// 	} finally {
// 		providersLoading.set(false);
// 	}
// }

// /**
//  * Set the active provider (used for creating new invoices)
//  */
// export function setActiveProvider(provider: ProviderData): void {
// 	console.log('[Provider Store] setActiveProvider called with:', provider.name);

// 	// Immediately save to localStorage to ensure persistence
// 	if (typeof window !== 'undefined') {
// 		localStorage.setItem('activeProvider', JSON.stringify(provider));
// 		console.log('[Provider Store] Saved to localStorage:', provider.name);
// 	}

// 	// Update the store - this should trigger all subscribers
// 	activeProvider.set(provider);
// 	console.log('[Provider Store] activeProvider.set() completed');

// 	// Force a verification
// 	if (typeof window !== 'undefined') {
// 		const stored = localStorage.getItem('activeProvider');
// 		console.log(
// 			'[Provider Store] Verified localStorage activeProvider:',
// 			stored ? JSON.parse(stored).name : 'null'
// 		);
// 	}
// }

// /**
//  * Clear the active provider
//  */
// export function clearActiveProvider(): void {
// 	activeProvider.set(null);

// 	// Immediately remove from localStorage
// 	if (typeof window !== 'undefined') {
// 		localStorage.removeItem('activeProvider');
// 		console.log('[Provider Store] Cleared activeProvider from localStorage');
// 	}
// }

// /**
//  * Reset providers to mock data (useful for development/debugging)
//  */
// export function resetProvidersToMock(): void {
// 	if (typeof window === 'undefined') return;

// 	console.log('Resetting providers to mock data...');
// 	localStorage.setItem('providers', JSON.stringify(mockProviders));
// 	providers.set(mockProviders);

// 	if (mockProviders.length > 0) {
// 		setActiveProvider(mockProviders[0]);
// 	}
// }

// /**
//  * Initialize providers on app startup
//  * Call this from your root layout
//  */
// export async function initializeProviders(): Promise<void> {
// 	await loadProviders();
// }
