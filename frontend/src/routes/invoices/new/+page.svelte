<script lang="ts">
	import InvoiceForm from '@/components/organisms/invoice-form/invoice-form.svelte';
	import { clients, invoices, providers } from '@/stores';
	import { api } from '@/services';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';
	import type { Invoice } from '@/types/invoice';
	import { page } from '$app/state';
	import { getDefaultDueDate, getTodayISOString } from '@/helpers';

	let isSaving = $state(false);
	let saveError = $state<string | null>(null);

	let duplicatedInvoice = $state<Invoice | undefined>(undefined);

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
			window.location.href = `/invoices/${createdInvoice.id}`;
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

	// Fetch data and check for duplication request
	async function loadData() {
		try {
			isLoading = true;
			const [clientsData, providersData] = await Promise.all([
				api.clients.getAllClients(fetch),
				api.providers.getAllProviders(fetch)
			]);

			clients.set(clientsData);
			providers.set(providersData);

			// check url for duplication_from parameter
			const duplicateId = page.url.searchParams.get('duplicate_from');
			if (duplicateId) {
				const sourceInvoice = await api.invoices.getInvoice(fetch, duplicateId);

				// prepare items: join description and details for the form editor
				// (the editor splits them on save, so must join here to show full content)
				const preparedItems = sourceInvoice.items.map((item) => ({
					...item,
					date: getTodayISOString(),
					description: item.description_detail
						? `${item.description}\n${item.description_detail}`
						: item.description
				}));

				// tramsform source invoice into a new draft
				duplicatedInvoice = {
					...sourceInvoice,
					id: '', // reset id to trigger 'new' mode
					date: getTodayISOString(), // reset issue date
					due: getDefaultDueDate(),
					status: 'draft',
					items: preparedItems,
					email_template_id: sourceInvoice.email_template_id || 'default'
				};
			}
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
		<InvoiceForm
			mode="create"
			invoice={duplicatedInvoice}
			{isSaving}
			onSave={handleSave}
			onCancel={handleCancel}
		/>
	{/if}
</div>
