<script lang="ts">
	import { page } from '$app/state';
	import InvoiceForm from '@/components/organisms/invoice-form/invoice-form.svelte';
	import { clients, providers } from '@/stores';
	import type { Invoice } from '@/types/invoice';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';
	import { api } from '@/services';

	interface Props {
		data: {
			invoice: Invoice | null;
			invErr?: string;
		};
	}

	let { data }: Props = $props();
	let invoice = $derived(data.invoice as Invoice);
	let error = $derived(data.invErr);

	function handleSave(data: any) {
		// TODO: Implement API call to update invoice
		console.log('Updating invoice:', invoice.id, data);
		// Redirect to invoice view after successful save
		window.history.back();
	}

	function handleCancel() {
		window.history.back();
	}

	// TODO: optimize fetching
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
	{#if error || !invoice}
		<ErrorAlert message={error || 'Invoice not found'} showBackButton={true} />
	{:else if loadError}
		<ErrorAlert message={loadError} showBackButton={true} />
	{:else if isLoading}
		<div class="flex items-center justify-center py-12">
			<p class="text-muted-foreground">Loading...</p>
		</div>
	{:else}
		<InvoiceForm mode="edit" {invoice} onSave={handleSave} onCancel={handleCancel} />
	{/if}
</div>
