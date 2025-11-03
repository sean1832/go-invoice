<script lang="ts">
	import InvoiceShelf from '@/components/organisms/shelf/invoice-shelf.svelte';
	import { invoices, filteredInvoices } from '$lib/stores/invoices';
	import { api } from '$lib/services';
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import Spinner from '@/components/atoms/spinner.svelte';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { Button } from '$lib/components/ui/button/index.js';

	let isLoading = $state(false);
	let errorMessage = $state<string | null>(null);
	let loadedInvoices = $state<any[]>([]);

	async function loadInvoices() {
		isLoading = true;
		errorMessage = null;

		try {
			const data = await api.invoices.getAllInvoices(fetch);
			invoices.set(data);
			loadedInvoices = data;
		} catch (error) {
			console.error('Failed to load invoices:', error);
			errorMessage =
				error instanceof Error ? error.message : 'Failed to load invoices. Please try again.';
		} finally {
			isLoading = false;
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
				{#if isLoading}
					<div class="flex flex-col items-center justify-center gap-4 py-12">
						<Spinner size={48} />
						<p class="text-muted-foreground">Loading invoices...</p>
					</div>
				{:else if errorMessage}
					<Alert.Root variant="destructive">
						<AlertCircleIcon />
						<Alert.Title>Error Loading Invoices</Alert.Title>
						<Alert.Description>{errorMessage}</Alert.Description>
					</Alert.Root>
					<div class="flex justify-center">
						<Button onclick={loadInvoices} class="w-full cursor-pointer">Retry</Button>
					</div>
				{:else}
					<InvoiceShelf data={loadedInvoices} />
				{/if}
			</div>
		</div>
	</div>
</div>
