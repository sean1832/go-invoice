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
		deletingInvoiceId?: string | null;
	}

	const { data, onError, onEdit, onDelete, deletingInvoiceId = null } = $props();

	let activeTab = $state<InvoiceStatus | 'all'>('all');

	const filteredByStatus = $derived(
		activeTab === 'all' ? data : data.filter((invoice: Invoice) => invoice.status === activeTab)
	);
	function createNewInvoice() {
		window.location.href = '/invoices/new';
	}
</script>

<div>
	<!-- Status tabs wrapper for invoices -->
	<Tabs.Root bind:value={activeTab}>
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

		<Tabs.Content value={activeTab} class="mt-4">
			<Shelf
				data={filteredByStatus}
				itemComponent={ShelfItem}
				searchFields={['id', 'client.name']}
				emptyMessage="No invoices found"
				searchPlaceholder="Search invoices by number or client..."
				{onError}
				{onEdit}
				{onDelete}
				deletingItemId={deletingInvoiceId}
			/>
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
