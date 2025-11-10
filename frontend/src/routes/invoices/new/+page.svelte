<script lang="ts">
	import InvoiceForm from '@/components/organisms/invoice-form/invoice-form.svelte';
	import { clients, invoices, providers } from '@/stores';
	import { api } from '@/services';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';
	import type { Invoice } from '@/types/invoice';

	let isSaving = $state(false);
	let saveError = $state<string | null>(null);

	async function handleSave(data: Invoice) {
		console.warn(data);
		try {
			isSaving = true;
			saveError = null;

			// Seperate description -> description_detail by first line break
			for (const item of data.items) {
				if (item.description.includes('\n')) {
					const lines = item.description.split('\n');
					item.description = lines[0];
					item.description_detail = lines.slice(1).join('\n');
				}
			}

			// Create invoice via API
			const createdInvoice = await api.invoices.createInvoice(fetch, data);

			// Update store with the new invoice
			invoices.update((current) => [...current, createdInvoice]);

			// Navigate back to invoice list
			window.history.back();
		} catch (err) {
			console.error('Failed to save invoice:', err);
			saveError = err instanceof Error ? err.message : 'Failed to save invoice. Please try again.';
		} finally {
			isSaving = false;
		}
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
		{#if saveError}
			<div class="mb-4">
				<ErrorAlert message={saveError} />
			</div>
		{/if}
		<InvoiceForm mode="create" {isSaving} onSave={handleSave} onCancel={handleCancel} />
	{/if}
</div>
