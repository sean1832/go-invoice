<script lang="ts">
	import Button from '@/components/ui/button/button.svelte';
	import * as Card from '@/components/ui/card';
	import type { ProviderData } from '@/types/invoice';
	import { activeProvider } from '@/stores';
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import { ProviderForm } from '@/components/organisms/profile-form';
	import { api } from '@/services';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';
	import { updateProvider } from '@/stores/provider';

	// Use reactive reference to activeProvider
	let currentProvider = $derived($activeProvider);

	let errorMessage = $state<string | null>(null);
	let isSaving = $state(false);

	async function handleSave(provider: ProviderData) {
		isSaving = true;
		try {
			await api.providers.updateProvider(fetch, provider.id, provider);

			// Update both the providers list AND the activeProvider
			updateProvider(provider);

			// Go back after saving
			window.history.back();
		} catch (err) {
			console.error('failed to update provider data: ', err);
			errorMessage = err instanceof Error ? err.message : 'Unknown error updating provider data.';
		} finally {
			isSaving = false;
		}
	}

	function handleCancel() {
		// Go back without saving
		window.history.back();
	}

	function navigateToNewProvider() {
		window.location.href = '/providers/new';
	}
</script>

<div class="container mx-auto max-w-4xl p-4">
	<div class="mb-6">
		<h1 class="text-3xl font-bold tracking-tight">Settings</h1>
		<p class="text-muted-foreground">Manage your provider profile and business information</p>
	</div>

	{#if errorMessage}
		<ErrorAlert
			message={errorMessage}
			title="Error"
			showRetryButton={false}
			onRetry={() => {
				errorMessage = null;
			}}
		/>
	{:else if currentProvider}
		<ProviderForm
			provider={currentProvider}
			mode="edit"
			onSave={handleSave}
			onCancel={handleCancel}
		/>
	{:else}
		<Card.Root>
			<Card.Content class="flex flex-col items-center justify-center py-12">
				<AlertCircleIcon class="mb-4 h-12 w-12 text-muted-foreground" />
				<h3 class="mb-2 text-lg font-semibold">No Provider Selected</h3>
				<p class="mb-6 text-center text-sm text-muted-foreground">
					You need to create a provider profile to start creating invoices
				</p>
				<Button onclick={navigateToNewProvider}>Create Provider Profile</Button>
			</Card.Content>
		</Card.Root>
	{/if}
</div>
