<script lang="ts">
	import InvoiceForm from '@/components/organisms/invoice-form/invoice-form.svelte';
	import { clients, providers } from '@/stores';
	import { api } from '@/services';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';

	function handleSave(data: any) {
		// TODO: Implement API call to save invoice
		console.log('Saving new invoice:', data);
		// Redirect to invoice list after successful save
		window.history.back();
	}

	function handleCancel() {
		window.history.back();
	}

	// // TODO: optimize fetching
	let isLoading = $state(true);
	let loadError = $state<string | null>(null);

	// Fetch both clients and providers data
	async function loadData() {
		try {
			isLoading = true;
			const [clientsData, providersData] = await Promise.all([
				api.clients.getAllClients(fetch),
				api.providers.getAllProviders(fetch)
			]);
			clients.set(clientsData);
			providers.set(providersData);
		} catch (err) {
			console.error('Failed to load data:', err);
			loadError = err instanceof Error ? err.message : 'Failed to load data. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	$effect(() => {
		loadData();
	});
</script>

<div class="container mx-auto max-w-5xl p-4">
	{#if loadError}
		<ErrorAlert message={loadError} showBackButton={true} />
	{:else if isLoading}
		<div class="flex items-center justify-center py-12">
			<p class="text-muted-foreground">Loading...</p>
		</div>
	{:else}
		<InvoiceForm mode="create" onSave={handleSave} onCancel={handleCancel} />
	{/if}
</div>
