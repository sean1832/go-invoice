<script lang="ts">
	import { ClientForm } from '@/components/organisms/profile-form';
	import type { ClientData } from '@/types/invoice';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';
	import { api } from '@/services';
	import { updateClient } from '@/stores/clients';

	interface Props {
		data: {
			client: ClientData | null;
			error?: string;
		};
	}

	let { data }: Props = $props();
	let client = $derived(data.client);
	let error = $derived(data.error);

	let isSaving = $state(false);
	let saveError = $state<string | null>(null);

	async function handleSave(client: ClientData) {
		console.log('Saving client:', client);
		isSaving = true;
		saveError = null;
		try {
			await api.clients.updateClient(fetch, client.id, client);

			// Update the store with the updated client
			updateClient(client);

			window.history.back();
		} catch (err) {
			console.error('failed to save client: ', err);
			saveError = err instanceof Error ? err.message : 'An unknown error occurred while saving.';
		} finally {
			isSaving = false;
		}
	}

	function handleCancel() {
		window.history.back();
	}
</script>

<div class="container mx-auto max-w-3xl p-4">
	{#if error || !client}
		<ErrorAlert message={error || 'Client not found'} showBackButton={true} />
	{:else}
		{#if saveError}
			<div class="mb-4">
				<ErrorAlert message={saveError} title="Save Failed" onRetry={() => (saveError = null)} />
			</div>
		{/if}
		<ClientForm
			{client}
			disable={isSaving}
			mode="edit"
			onSave={handleSave}
			onCancel={handleCancel}
		/>
	{/if}
</div>
