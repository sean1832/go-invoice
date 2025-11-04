<script lang="ts">
	import * as Popover from '@/components/ui/popover';
	import Button from '@/components/ui/button/button.svelte';
	import Separator from '@/components/ui/separator/separator.svelte';
	import UserIcon from '@lucide/svelte/icons/circle-user-round';
	import CheckIcon from '@lucide/svelte/icons/check';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import SettingsIcon from '@lucide/svelte/icons/settings';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import type { ProviderData } from '@/types/invoice';
	import { providers as providersStore, activeProvider, providers } from '@/stores';
	import { get } from 'svelte/store';
	import { onMount } from 'svelte';
	import { api } from '@/services';

	let open = $state(false);
	let currentProvider = $state<ProviderData | null>(null);
	let currentProviders = $state<ProviderData[]>([]);
	let updateKey = $state(0); // Force re-render key

	// Sync with stores using $effect - this ensures reactivity
	$effect(() => {
		const unsubProvider = activeProvider.subscribe((value) => {
			currentProvider = value;
		});
		return unsubProvider;
	});

	$effect(() => {
		const unsubProviders = providersStore.subscribe((value) => {
			currentProviders = value;
		});
		return unsubProviders;
	});

	let isLoading = $state(true);
	let loadError = $state<string | null>(null);

	async function loadProviders() {
		try {
			isLoading = true;
			const data = await api.providers.getAllProviders(fetch);
			providers.set(data);
		} catch (err) {
			console.error(`[profile selector] failed to load data: `, err);
			loadError = err instanceof Error ? err.message : 'Failed to load data. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	onMount(async () => {
		await loadProviders();

		// Check if there's an active provider in localStorage
		const storedActiveProvider = localStorage.getItem('activeProvider');

		// Get current value from store using get()
		const currentStoreValue = get(activeProvider);

		// If we have a stored active provider, refresh it from backend to get latest data
		if (storedActiveProvider || currentStoreValue) {
			let providerId: string = '';
			if (currentStoreValue && currentStoreValue.id) {
				providerId = currentStoreValue.id;
			} else if (storedActiveProvider) {
				providerId = (JSON.parse(storedActiveProvider) as ProviderData).id;
			}

			if (providerId) {
				try {
					const freshProvider = await api.providers.getProvider(fetch, providerId);
					activeProvider.set(freshProvider);
				} catch (err) {
					console.error('[provider selector] failed to refresh active provider:', err);
					// Fall back to using stored data if fetch fails
				}
			}
		} else if (currentProviders.length > 0) {
			// Only set default if nothing is stored AND store is empty AND we have providers
			activeProvider.set(currentProviders[0]); // default to first one
		}
	});

	function handleProviderSelect(provider: ProviderData) {
		activeProvider.set(provider);
		currentProvider = provider; // Force immediate update
		updateKey++; // Force re-render
		// Give the state time to update before closing
		setTimeout(() => {
			open = false;
		}, 0);
	}

	async function handleDeleteProvider(provider: ProviderData, event: MouseEvent) {
		event.stopPropagation(); // Prevent selecting the provider

		if (confirm(`Are you sure you want to delete ${provider.name}?`)) {
			try {
				await api.providers.deleteProvider(fetch, provider.id);

				// Remove from store - UI updates automatically
				providersStore.update((items) => items.filter((p) => p.id !== provider.id));

				// If this was the active provider, clear it or set another one
				const currentProviderValue = get(activeProvider);
				if (currentProviderValue?.id === provider.id) {
					const updatedProviders = get(providersStore);
					if (updatedProviders.length > 0) {
						activeProvider.set(updatedProviders[0]);
					} else {
						localStorage.removeItem('activeProvider');
						window.location.href = '/providers/new';
					}
				}
			} catch (err) {
				console.error('Failed to delete provider:', err);
				alert(err instanceof Error ? err.message : `Failed to delete provider ${provider.name}`);
			}
		}
	}

	function navigateToSettings() {
		open = false;
		window.location.href = '/settings';
	}

	function navigateToNewProvider() {
		open = false;
		window.location.href = '/providers/new';
	}
</script>

{#key updateKey}
	<Popover.Root bind:open>
		<Popover.Trigger
			class="inline-flex h-9 w-9 cursor-pointer 
			items-center justify-center rounded-full border 
			border-input bg-background text-sm font-medium 
			whitespace-nowrap shadow-sm ring-offset-background transition-colors 
			hover:bg-accent hover:text-accent-foreground focus-visible:ring-2 
			focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:outline-none 
			disabled:pointer-events-none disabled:opacity-50"
		>
			<UserIcon class="h-6 w-6" />
		</Popover.Trigger>
		<Popover.Content class="w-80 p-0" align="end">
			<div class="flex flex-col">
				<!-- Current Provider Header -->
				<div class="px-4 py-3">
					{#if currentProvider}
						<div class="flex flex-col">
							<span class="font-semibold">{currentProvider.name}</span>
							<span class="text-xs text-muted-foreground">{currentProvider.email}</span>
						</div>
					{:else}
						<span class="text-sm text-muted-foreground">No Provider Selected</span>
					{/if}
				</div>

				<Separator />

				<!-- Provider List -->
				{#if currentProviders.length > 0}
					<div class="max-h-64 overflow-y-auto">
						<div class="px-2 py-2">
							<p class="mb-2 px-2 text-xs font-medium text-muted-foreground">Switch Provider</p>
							{#each currentProviders as provider (provider.id)}
								<div class="group relative">
									<button
										onclick={() => handleProviderSelect(provider)}
										class="relative z-0 flex w-full items-center justify-between rounded-sm px-2 py-2 pr-10 text-left text-sm transition-colors hover:bg-accent hover:text-accent-foreground"
									>
										<div class="flex flex-col">
											<span class="font-medium">{provider.name}</span>
											<span class="text-xs text-muted-foreground">{provider.email}</span>
										</div>
										{#if currentProvider?.id === provider.id}
											<CheckIcon class="h-4 w-4 text-primary" />
										{/if}
									</button>
									<button
										onclick={(e) => handleDeleteProvider(provider, e)}
										class="hover:text-destructive-foreground absolute top-1/2 right-2 z-10 -translate-y-1/2 rounded-sm p-1 opacity-0 transition-opacity group-hover:opacity-100 hover:bg-destructive"
										title="Delete provider"
									>
										<TrashIcon class="h-4 w-4" />
									</button>
								</div>
							{/each}
						</div>
					</div>
					<Separator />
				{/if}

				<!-- Actions -->
				<div class="px-2 py-2">
					<button
						onclick={navigateToSettings}
						class="flex w-full items-center rounded-sm px-2 py-2 text-sm transition-colors hover:bg-accent hover:text-accent-foreground"
					>
						<SettingsIcon class="mr-2 h-4 w-4" />
						<span>Settings</span>
					</button>
					<button
						onclick={navigateToNewProvider}
						class="flex w-full items-center rounded-sm px-2 py-2 text-sm transition-colors hover:bg-accent hover:text-accent-foreground"
					>
						<PlusIcon class="mr-2 h-4 w-4" />
						<span>Add Provider</span>
					</button>
				</div>
			</div>
		</Popover.Content>
	</Popover.Root>
{/key}
