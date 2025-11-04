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

export function removeProvider(id: string): void {
	providers.update((items) => items.filter((inv) => inv.id !== id));
}

export function updateProvider(item: ProviderData): void {
	providers.update((items) => {
		return items.map((inv) => (inv.id === item.id ? { ...inv, ...item } : inv));
	});
	activeProvider.set(item);
}
