<script lang="ts">
	import Shelf from '@/components/organisms/shelf/shelf.svelte';
	import ShelfItem from '@/components/organisms/shelf/invoice-item.svelte';
	import * as Tabs from '@/components/ui/tabs';
	import type { Invoice, InvoiceStatus } from '@/types/invoice';
	import Button from '@/components/ui/button/button.svelte';
	import PlusIcon from '@lucide/svelte/icons/plus';

	interface Props {
		data: Invoice[];
		onError?: (message: string) => void;
		onEdit: (item: Invoice) => void;
		onDelete: (item: Invoice) => void;
		onDownload?: (item: Invoice) => void;
		deletingInvoiceId?: string | null;
		isDownloading: boolean;
		// Pagination props
		page?: number;
		totalPages?: number;
		totalCount?: number;
		onPageChange?: (page: number) => void;
		// Status filter props (controlled by parent for server-side filtering)
		activeStatus?: InvoiceStatus | 'all';
		onStatusChange?: (status: InvoiceStatus | 'all') => void;
	}

	const {
		data,
		onError,
		onEdit,
		onDuplicate,
		onDelete,
		onDownload,
		deletingInvoiceId = null,
		isDownloading,
		page = 1,
		totalPages = 1,
		totalCount = 0,
		onPageChange,
		activeStatus = 'all',
		onStatusChange
	} = $props();

	// Show pagination only if there's more than one page
	const showPagination = $derived(totalPages > 1);

	function createNewInvoice() {
		window.location.href = '/invoices/new';
	}

	function handlePageChange(newPage: number) {
		if (onPageChange && newPage >= 1 && newPage <= totalPages) {
			onPageChange(newPage);
		}
	}

	function handleTabChange(value: string) {
		if (onStatusChange) {
			onStatusChange(value as InvoiceStatus | 'all');
		}
	}
</script>

<div>
	<!-- Status tabs wrapper for invoices -->
	<Tabs.Root value={activeStatus} onValueChange={handleTabChange}>
		<div class="flex items-center justify-between gap-4">
			<Tabs.List class="grid flex-1 grid-cols-3 md:w-auto md:flex-initial">
				<Tabs.Trigger value="all">All</Tabs.Trigger>
				<Tabs.Trigger value="draft">Drafts</Tabs.Trigger>
				<Tabs.Trigger value="send">Sent</Tabs.Trigger>
			</Tabs.List>

			<!-- Desktop Create Button -->
			<Button onclick={createNewInvoice} class="hidden md:flex" size="default">
				<PlusIcon class="mr-2 h-4 w-4" />
				New Invoice
			</Button>
		</div>

		<Tabs.Content value={activeStatus} class="mt-4">
			<Shelf
				{data}
				searchFields={['id', 'client.name']}
				emptyMessage="No invoices found"
				searchPlaceholder="Search invoices by number or client..."
				deletingItemId={deletingInvoiceId}
			>
				{#snippet children(item: Invoice, isDeleting: boolean)}
					<ShelfItem
						{item}
						{onError}
						{onEdit}
						{onDuplicate}
						{onDelete}
						{onDownload}
						{isDeleting}
						{isDownloading}
					/>
				{/snippet}
			</Shelf>

			<!-- Pagination Controls -->
			{#if showPagination}
				<div class="mt-6 flex flex-col items-center gap-2">
					<div class="flex items-center gap-1">
						<Button
							variant="outline"
							size="icon"
							onclick={() => handlePageChange(page - 1)}
							disabled={page <= 1}
							aria-label="Previous page"
						>
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m15 18-6-6 6-6"/></svg>
						</Button>

						{#each { length: totalPages } as _, i}
							{@const pageNum = i + 1}
							{#if pageNum === 1 || pageNum === totalPages || (pageNum >= page - 1 && pageNum <= page + 1)}
								<Button
									variant={pageNum === page ? 'default' : 'outline'}
									size="icon"
									onclick={() => handlePageChange(pageNum)}
									aria-current={pageNum === page ? 'page' : undefined}
								>
									{pageNum}
								</Button>
							{:else if pageNum === page - 2 || pageNum === page + 2}
								<span class="text-muted-foreground px-2">...</span>
							{/if}
						{/each}

						<Button
							variant="outline"
							size="icon"
							onclick={() => handlePageChange(page + 1)}
							disabled={page >= totalPages}
							aria-label="Next page"
						>
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m9 18 6-6-6-6"/></svg>
						</Button>
					</div>
					<p class="text-muted-foreground text-sm">
						Showing page {page} of {totalPages} ({totalCount} invoices)
					</p>
				</div>
			{/if}
		</Tabs.Content>
	</Tabs.Root>

	<!-- Mobile Floating Action Button -->
	<Button
		onclick={createNewInvoice}
		size="lg"
		class="fixed right-8 bottom-8 h-14 w-14 rounded-full shadow-lg transition-shadow hover:shadow-xl md:hidden"
		title="Create new invoice"
	>
		<PlusIcon class="h-6 w-6" />
	</Button>
</div>


