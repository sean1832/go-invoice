<script lang="ts">
	import ErrorAlert from '@/components/molecules/error-alert.svelte';
	import { ProviderForm } from '@/components/organisms/profile-form';
	import { api } from '@/services';
	import { addProvider } from '@/stores/provider';
	import type { ProviderData } from '@/types/invoice';

	let errorMessage = $state<string | null>(null);
	let isSaving = $state(false);

	function generateProviderId(provider: ProviderData): string {
		const id = provider.name
			.toLowerCase()
			.replace(/\s+/g, '_')
			.replace(/[^a-z0-9_]/g, ''); // Remove invalid characters
		return id;
	}

	async function handleSave(provider: ProviderData) {
		isSaving = true;
		try {
			// Generate ID if creating new provider
			if (!provider.id || provider.id.trim() === '') {
				provider.id = generateProviderId(provider);
			}
			await api.providers.createProvider(fetch, provider);

			// Add to providers list AND set as activeProvider
			addProvider(provider);

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
		window.history.back();
	}
</script>

<div class="container mx-auto max-w-3xl p-4">
	{#if errorMessage}
		<ErrorAlert
			message={errorMessage}
			title="Error"
			showRetryButton={false}
			onRetry={() => {
				errorMessage = null;
			}}
		/>
	{/if}
	<ProviderForm mode="create" onSave={handleSave} onCancel={handleCancel} />
</div>
