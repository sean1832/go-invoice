<script lang="ts">
	import { page } from '$app/state';
	import InvoiceForm from '@/components/organisms/invoice-form/invoice-form.svelte';
	import { clients, providers, invoices } from '@/stores';
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

	let isSaving = $state(false);
	let saveError = $state<string | null>(null);

	async function handleSave(updatedData: Invoice) {
		try {
			isSaving = true;
			saveError = null;

			// Update invoice via API
			const updatedInvoice = await api.invoices.updateInvoice(fetch, invoice.id, updatedData);

			// Update store with the updated invoice
			invoices.update((current) =>
				current.map((inv) => (inv.id === updatedInvoice.id ? updatedInvoice : inv))
			);

			// Navigate back to invoice view
			window.history.back();
		} catch (err) {
			console.error('Failed to update invoice:', err);
			saveError =
				err instanceof Error ? err.message : 'Failed to update invoice. Please try again.';
		} finally {
			isSaving = false;
		}
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
		{#if saveError}
			<div class="mb-4">
				<ErrorAlert message={saveError} />
			</div>
		{/if}
		<InvoiceForm mode="edit" {invoice} {isSaving} onSave={handleSave} onCancel={handleCancel} />
	{/if}
</div>
