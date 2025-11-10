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
// NOTE: This cached data should be validated against server on component mount
// Components should always fetch fresh data and validate cached provider still exists
if (typeof window !== 'undefined') {
	const storedActiveProvider = localStorage.getItem('activeProvider');
	if (storedActiveProvider) {
		try {
			const parsed = JSON.parse(storedActiveProvider);
			// Only set if we have a valid ID - component will validate against server
			if (parsed && parsed.id) {
				activeProvider.set(parsed);
			}
		} catch (error) {
			console.error('Failed to parse active provider from localStorage:', error);
			// Clear corrupted cache
			localStorage.removeItem('activeProvider');
		}
	}
}

// Subscribe to save active provider changes to localStorage
// This keeps cache in sync with user's selections
if (typeof window !== 'undefined') {
	activeProvider.subscribe((value) => {
		if (value && value.id) {
			localStorage.setItem('activeProvider', JSON.stringify(value));
		} else {
			localStorage.removeItem('activeProvider');
		}
	});
}

export function removeProvider(id: string): void {
	providers.update((items) => {
		const itemsArray = Array.isArray(items) ? items : [];
		return itemsArray.filter((inv) => inv.id !== id);
	});
}

export function addProvider(item: ProviderData): void {
	providers.update((items) => {
		// Ensure items is always an array
		const itemsArray = Array.isArray(items) ? items : [];

		// Check if provider already exists
		const exists = itemsArray.some((p) => p.id === item.id);
		if (exists) {
			// Update existing
			return itemsArray.map((p) => (p.id === item.id ? { ...p, ...item } : p));
		} else {
			// Add new
			return [...itemsArray, item];
		}
	});
	activeProvider.set(item);
}

export function updateProvider(item: ProviderData): void {
	providers.update((items) => {
		const itemsArray = Array.isArray(items) ? items : [];
		return itemsArray.map((inv) => (inv.id === item.id ? { ...inv, ...item } : inv));
	});
	activeProvider.set(item);
}
