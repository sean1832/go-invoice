import { writable } from 'svelte/store';
import type { ProviderData } from '@/types/invoice';

// Load active provider from localStorage
function loadActiveProvider(): ProviderData | null {
	if (typeof window === 'undefined') return null;

	const stored = localStorage.getItem('activeProvider');
	return stored ? JSON.parse(stored) : null;
}

// Create the store
export const activeProvider = writable<ProviderData | null>(loadActiveProvider());

// Subscribe to save changes to localStorage
if (typeof window !== 'undefined') {
	activeProvider.subscribe((value) => {
		if (value) {
			localStorage.setItem('activeProvider', JSON.stringify(value));
		} else {
			localStorage.removeItem('activeProvider');
		}
	});
}

// Helper functions
export function setActiveProvider(provider: ProviderData) {
	activeProvider.set(provider);
}

export function clearActiveProvider() {
	activeProvider.set(null);
}

// Load available providers (mock for now, will be from API later)
export async function loadProviders(): Promise<ProviderData[]> {
	// TODO: Replace with actual API call
	// For now, return mock data from localStorage or empty array
	if (typeof window === 'undefined') return [];

	const stored = localStorage.getItem('providers');
	return stored ? JSON.parse(stored) : [];
}

// Save provider to available providers list
export async function saveProvider(provider: ProviderData): Promise<void> {
	if (typeof window === 'undefined') return;

	const providers = await loadProviders();
	const index = providers.findIndex((p) => p.id === provider.id);

	if (index >= 0) {
		providers[index] = provider;
	} else {
		providers.push(provider);
	}

	localStorage.setItem('providers', JSON.stringify(providers));

	// If this is the active provider or no active provider exists, update it
	const current = loadActiveProvider();
	if (!current || current.id === provider.id) {
		setActiveProvider(provider);
	}
}

// Initialize with mock data if no providers exist
export function initializeMockProviders(): void {
	if (typeof window === 'undefined') return;

	const providers = localStorage.getItem('providers');
	if (!providers) {
		const mockProviders: ProviderData[] = [];

		localStorage.setItem('providers', JSON.stringify(mockProviders));

		// Set the first one as active if no active provider
		const current = loadActiveProvider();
		if (!current) {
			setActiveProvider(mockProviders[0]);
		}
	}
}
