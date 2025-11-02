<script lang="ts">
	import Shelf from '@/components/custom/shelf/shelf.svelte';
	import ShelfItem from '@/components/custom/shelf/shelf-item.svelte';
	import * as Tabs from '@/components/ui/tabs';
	import type { Invoice, InvoiceStatus } from '@/types/invoice';

	interface Props {
		data: Invoice[];
	}

	const { data } = $props();

	let activeTab = $state<InvoiceStatus | 'all'>('all');

	const filteredByStatus = $derived(
		activeTab === 'all' ? data : data.filter((invoice: Invoice) => invoice.status === activeTab)
	);
</script>

<!-- Status tabs wrapper for invoices -->
<Tabs.Root bind:value={activeTab}>
	<Tabs.List class="grid w-full grid-cols-3">
		<Tabs.Trigger value="all">All</Tabs.Trigger>
		<Tabs.Trigger value="draft">Drafts</Tabs.Trigger>
		<Tabs.Trigger value="send">Sent</Tabs.Trigger>
	</Tabs.List>

	<Tabs.Content value={activeTab} class="mt-4">
		<Shelf
			data={filteredByStatus}
			itemComponent={ShelfItem}
			searchFields={['id', 'client.name']}
			emptyMessage="No invoices found"
			searchPlaceholder="Search invoices by number or client..."
		/>
	</Tabs.Content>
</Tabs.Root>
