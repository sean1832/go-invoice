<script lang="ts">
	import InvoiceShelf from '@/components/organisms/shelf/invoice-shelf.svelte';
	import { invoices, filteredInvoices, removeInvoice } from '$lib/stores/invoices';
	import { api } from '$lib/services';
	import Spinner from '@/components/atoms/spinner.svelte';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';
	import type { Invoice } from '@/types/invoice';

	let isLoading = $state(false);
	let errorMessage = $state<string | null>(null);

	async function loadInvoices() {
		isLoading = true;
		errorMessage = null;

		try {
			const data = await api.invoices.getAllInvoices(fetch);
			invoices.set(data);
		} catch (error) {
			console.error('Failed to load invoices:', error);
			errorMessage =
				error instanceof Error ? error.message : 'Failed to load invoices. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	function handleError(message: string) {
		errorMessage = message;
	}

	function clearError() {
		errorMessage = null;
	}

	let deletingInvoiceId = $state<string | null>(null);
	async function handleDelete(invoice: Invoice) {
		deletingInvoiceId = invoice.id;

		try {
			await api.invoices.deleteInvoice(fetch, invoice.id);
			// Remove from store - UI updates automatically
			removeInvoice(invoice.id);
		} catch (err) {
			const message = err instanceof Error ? err.message : `Failed to delete invoice ${invoice.id}`;
			console.error('Failed to delete invoice:', err);
			handleError(message);
		} finally {
			deletingInvoiceId = null;
		}
	}

	async function handleEdit(invoice: Invoice) {
		window.location.href = `/invoices/${invoice.id}/edit`;
	}

	async function handleDownload(invoice: Invoice) {
		try {
			const blob = await api.invoices.downloadPdf(fetch, invoice.id);

			// Create a temporary link to trigger the download
			const link = document.createElement('a');
			const url = URL.createObjectURL(blob);

			link.href = url;
			link.download = `${invoice.id}.pdf`;

			document.body.appendChild(link);
			link.click();

			// Clean up
			document.body.removeChild(link);
			URL.revokeObjectURL(url);
		} catch (error) {
			console.error('Error downloading invoice PDF:', error);
			handleError(
				error instanceof Error ? error.message : 'Failed to download PDF. Please try again.'
			);
		}
	}
	// Load invoices on mount
	$effect(() => {
		loadInvoices();
	});
</script>

<div class="p-4">
	<div class="container mx-auto flex justify-center">
		<div class="w-full max-w-4xl">
			<div class="mb-6">
				<h1 class="text-3xl font-bold tracking-tight">Invoices</h1>
				<p class="text-muted-foreground">Manage and track your invoices</p>
			</div>
			<div class="flex flex-col gap-4">
				{#if errorMessage}
					<ErrorAlert
						message={errorMessage}
						title="Error"
						showRetryButton={false}
						onRetry={clearError}
					/>
				{/if}

				{#if isLoading}
					<div class="flex flex-col items-center justify-center gap-4 py-12">
						<Spinner size={48} />
						<p class="text-muted-foreground">Loading invoices...</p>
					</div>
				{:else}
					<InvoiceShelf
						data={$filteredInvoices}
						onError={handleError}
						onDelete={handleDelete}
						onEdit={handleEdit}
						onDownload={handleDownload}
						{deletingInvoiceId}
					/>
				{/if}
			</div>
		</div>
	</div>
</div>
