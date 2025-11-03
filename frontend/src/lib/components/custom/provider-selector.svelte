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
	import {
		providers as providersStore,
		activeProvider,
		setActiveProvider,
		loadProviders
	} from '@/stores';
	import { onMount } from 'svelte';

	let open = $state(false);

	// Use reactive references to stores
	let currentProvider = $derived($activeProvider);
	let providers = $derived($providersStore);

	// Reload providers when popover opens
	$effect(() => {
		if (open) {
			console.log('Reloading providers...');
			loadProviders();
		}
	});

	onMount(async () => {
		await loadProviders();
	});

	function handleProviderSelect(provider: ProviderData) {
		setActiveProvider(provider);
		open = false;
	}

	async function handleDeleteProvider(provider: ProviderData, event: MouseEvent) {
		event.stopPropagation(); // Prevent selecting the provider

		if (confirm(`Are you sure you want to delete ${provider.name}?`)) {
			// Remove from localStorage
			const updatedProviders = providers.filter((p) => p.id !== provider.id);
			localStorage.setItem('providers', JSON.stringify(updatedProviders));

			// If this was the active provider, clear it or set another one
			if (currentProvider?.id === provider.id) {
				if (updatedProviders.length > 0) {
					setActiveProvider(updatedProviders[0]);
				} else {
					localStorage.removeItem('activeProvider');
					window.location.href = '/providers/new';
				}
			}

			// Reload the list
			providers = await loadProviders();
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

<Popover.Root bind:open>
	<Popover.Trigger>
		<Button class="rounded-full p-1" variant="outline">
			<UserIcon class="h-6 w-6" />
		</Button>
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
			{#if providers.length > 0}
				<div class="max-h-64 overflow-y-auto">
					<div class="px-2 py-2">
						<p class="mb-2 px-2 text-xs font-medium text-muted-foreground">Switch Provider</p>
						{#each providers as provider (provider.id)}
							<div class="group relative">
								<button
									onclick={() => handleProviderSelect(provider)}
									class="flex w-full items-center justify-between rounded-sm px-2 py-2 pr-10 text-left text-sm transition-colors hover:bg-accent hover:text-accent-foreground"
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
									class="hover:text-destructive-foreground absolute top-1/2 right-2 -translate-y-1/2 rounded-sm p-1 opacity-0 transition-opacity group-hover:opacity-100 hover:bg-destructive"
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
