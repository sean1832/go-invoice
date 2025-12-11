<script lang="ts">
	import * as Item from '$lib/components/ui/item/index.js';
	import { Badge } from '@/components/ui/badge';
	import Button from '@/components/ui/button/button.svelte';
	import { ConfirmDialog } from '@/components/molecules';
	import type { Invoice } from '@/types/invoice';
	import { CopyIcon, TrashIcon, PencilIcon } from '@lucide/svelte/icons';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import Spinner from '@/components/atoms/spinner.svelte';

	interface Props {
		item: Invoice;
		onError?: (message: string) => void;
		onEdit: (item: Invoice) => void;
		onDuplicate: (item: Invoice) => void;
		onDelete: (item: Invoice) => void;
		onDownload?: (item: Invoice) => void;
		isDeleting?: boolean;
		isDownloading: boolean;
	}

	let {
		item: invoice,
		onError,
		onEdit,
		onDuplicate,
		onDelete,
		onDownload,
		isDeleting = false,
		isDownloading
	}: Props = $props();

	let deleteDialogOpen = $state(false);

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

	function duplicateInvoice() {
		onDuplicate(invoice);
	}

	function openDeleteDialog() {
		deleteDialogOpen = true;
	}

	function confirmDelete() {
		onDelete(invoice);
	}

	function downloadInvoice() {
		if (onDownload) {
			onDownload(invoice);
		}
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
			<PencilIcon class="h-4 w-4" />
		</Button>
		<Button
			variant="ghost"
			size="sm"
			onclick={duplicateInvoice}
			title="duplicate"
			disabled={isDeleting}
		>
			<CopyIcon class="h-4 w-4" />
		</Button>
		<Button
			variant="ghost"
			size="sm"
			onclick={downloadInvoice}
			title="Download"
			disabled={isDeleting || isDownloading}
		>
			{#if isDownloading}
				<Spinner class="h-4 w-4" />
			{:else}
				<DownloadIcon class="h-4 w-4" />
			{/if}
		</Button>
		<Button
			variant="ghost"
			size="sm"
			onclick={openDeleteDialog}
			title="Delete"
			disabled={isDeleting}
		>
			{#if isDeleting}
				<span class="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent"
				></span>
			{:else}
				<TrashIcon class="h-4 w-4 text-destructive" />
			{/if}
		</Button>
	</Item.Actions>
</Item.Root>

<ConfirmDialog
	bind:open={deleteDialogOpen}
	title="Delete Invoice"
	description="Are you sure you want to delete invoice <strong>{invoice.id}</strong>? This action cannot be undone."
	confirmText="Delete"
	cancelText="Cancel"
	confirmVariant="destructive"
	onConfirm={confirmDelete}
/>
