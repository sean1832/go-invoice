<script lang="ts">
	import { Badge } from '@/components/ui/badge';
	import Button from '@/components/ui/button/button.svelte';
	import InvoiceDisplayCard from '@/components/organisms/invoice-display/invoice-display-card.svelte';
	import type { Invoice } from '@/types/invoice';
	import EditIcon from '@lucide/svelte/icons/pencil';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import CopyIcon from '@lucide/svelte/icons/copy';
	import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';

	interface Props {
		data: {
			invoice: Invoice | null;
			error?: string;
		};
	}

	let { data }: Props = $props();
	let invoice = $derived(data.invoice as Invoice);
	let error = $derived(data.error);

	// Helper functions
	function getStatusVariant(status: Invoice['status']): 'default' | 'secondary' {
		return status === 'send' ? 'default' : 'secondary';
	}

	function getStatusLabel(status: Invoice['status']): string {
		return status.charAt(0).toUpperCase() + status.slice(1);
	}

	// Action handlers
	function goBack() {
		// window.history.back();
		window.location.href = '/';
	}

	function editInvoice() {
		window.location.href = `/invoices/${invoice.id}/edit`;
	}

	function downloadInvoice() {
		// TODO: Implement download logic
		console.log('Download invoice:', invoice.id);
	}

	function duplicateInvoice() {
		// TODO: Implement duplicate logic
		console.log('Duplicate invoice:', invoice.id);
	}
</script>

<div class="container mx-auto max-w-5xl p-4">
	{#if error || !invoice}
		<!-- Error State -->
		<div class="mb-6">
			<Button variant="ghost" size="sm" onclick={goBack} class="mb-4">
				<ArrowLeftIcon class="mr-2 h-4 w-4" />
				Back
			</Button>
			<Alert.Root variant="destructive">
				<AlertCircleIcon class="h-4 w-4" />
				<Alert.Title>Error</Alert.Title>
				<Alert.Description>
					{error || 'Invoice not found'}
				</Alert.Description>
			</Alert.Root>
		</div>
	{:else}
		<!-- Action Buttons Bar -->
		<div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
			<Button variant="ghost" size="sm" onclick={goBack} class="self-start">
				<ArrowLeftIcon class="mr-2 h-4 w-4" />
				Back
			</Button>
			<div class="flex flex-wrap gap-2">
				<Badge variant={getStatusVariant(invoice.status)} class="px-3 py-1 text-sm">
					{getStatusLabel(invoice.status)}
				</Badge>
				<Button variant="default" size="sm" onclick={editInvoice}>
					<EditIcon class="mr-2 h-4 w-4" />
					<span class="hidden sm:inline">Edit</span>
				</Button>
				<Button variant="outline" size="sm" onclick={downloadInvoice}>
					<DownloadIcon class="h-4 w-4 sm:mr-2" />
					<span class="hidden sm:inline">Download PDF</span>
				</Button>
				<Button variant="outline" size="sm" onclick={duplicateInvoice}>
					<CopyIcon class="h-4 w-4 sm:mr-2" />
					<span class="hidden sm:inline">Duplicate</span>
				</Button>
			</div>
		</div>

		<!-- Invoice Display Component -->
		<InvoiceDisplayCard {invoice} />
	{/if}
</div>
