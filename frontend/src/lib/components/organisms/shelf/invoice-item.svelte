<script lang="ts">
	import * as Item from '$lib/components/ui/item/index.js';
	import { Badge } from '@/components/ui/badge';
	import Button from '@/components/ui/button/button.svelte';
	import type { Invoice } from '@/types/invoice';
	import EyeIcon from '@lucide/svelte/icons/eye';
	import EditIcon from '@lucide/svelte/icons/pencil';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import CopyIcon from '@lucide/svelte/icons/copy';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import { api } from '@/services';
	import { removeInvoice } from '@/stores/invoices';

	interface Props {
		item: Invoice;
		onError?: (message: string) => void;
		onEdit: (item: Invoice) => void;
		onDelete: (item: Invoice) => void;
		isDeleting?: boolean;
	}

	let { item: invoice, onError, onEdit, onDelete, isDeleting = false }: Props = $props();

	// Helper function to format currency
	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-AU', {
			style: 'currency',
			currency: 'AUD'
		}).format(amount);
	}

	// Helper function to format date
	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

	// Get badge variant based on status
	function getStatusVariant(status: Invoice['status']): 'default' | 'secondary' {
		return status === 'send' ? 'default' : 'secondary';
	}

	// Navigation handlers (to be implemented later)
	function viewInvoice() {
		window.location.href = `/invoices/${invoice.id}`;
	}

	function editInvoice() {
		onEdit(invoice);
	}

	async function deleteInvoice() {
		if (confirm(`Are you sure you want to delete ${invoice.id}?`)) {
			onDelete(invoice);
		}
		return;
	}

	function downloadInvoice() {
		// TODO: Implement download logic
		console.log('Download invoice:', invoice.id);
	}
</script>

<Item.Root variant="outline" class="cursor-pointer transition-colors hover:bg-muted/50">
	<Item.Content class="flex-1" onclick={viewInvoice}>
		<div class="flex items-start justify-between gap-4">
			<div class="flex-1">
				<div class="mb-1 flex items-center gap-2">
					<Item.Title class="text-lg font-semibold">
						{invoice.id}
					</Item.Title>
					<Badge variant={getStatusVariant(invoice.status)}>
						{invoice.status.charAt(0).toUpperCase() + invoice.status.slice(1)}
					</Badge>
				</div>
				<div class="text-sm text-muted-foreground">
					<div class="flex flex-col gap-1">
						<div class="flex items-center gap-2">
							<span class="font-medium">{invoice.client.name}</span>
						</div>
						<div class="flex items-center gap-4 text-xs">
							<span>Date: {formatDate(invoice.date)}</span>
							<span>Due: {formatDate(invoice.due)}</span>
						</div>
					</div>
				</div>
			</div>

			<div class="flex flex-col items-end gap-1">
				<div class="text-xl font-bold">
					{formatCurrency(invoice.pricing.total)}
				</div>
				<div class="text-xs text-muted-foreground">
					Subtotal: {formatCurrency(invoice.pricing.subtotal)}
				</div>
			</div>
		</div>
	</Item.Content>

	<Item.Actions class="flex gap-1">
		<Button variant="ghost" size="sm" onclick={editInvoice} title="Edit" disabled={isDeleting}>
			<EditIcon class="h-4 w-4" />
		</Button>
		<Button
			variant="ghost"
			size="sm"
			onclick={downloadInvoice}
			title="Download"
			disabled={isDeleting}
		>
			<DownloadIcon class="h-4 w-4" />
		</Button>
		<Button variant="ghost" size="sm" onclick={deleteInvoice} title="Delete" disabled={isDeleting}>
			{#if isDeleting}
				<span class="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent"
				></span>
			{:else}
				<TrashIcon class="h-4 w-4 text-destructive" />
			{/if}
		</Button>
	</Item.Actions>
</Item.Root>
