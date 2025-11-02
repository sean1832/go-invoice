<script lang="ts">
	import ShelfItem from '@/components/custom/shelf/shelf-item.svelte';
	import * as Tabs from '@/components/ui/tabs';
	import Input from '@/components/ui/input/input.svelte';
	import type { Invoice, InvoiceStatus } from '@/types/invoice';
	import SearchIcon from '@lucide/svelte/icons/search';

	// Mock data - this will be replaced with real API calls later
	const mockInvoices: Invoice[] = [
		{
			id: 'INV-251103001',
			status: 'draft',
			date: '2025-11-01',
			due: '2025-11-30',
			client: {
				name: 'Acme Corporation',
				email: 'contact@acme.com',
				address: '123 Business St'
			},
			provider: {
				name: 'Your Company',
				email: 'you@company.com'
			},
			items: [],
			pricing: {
				subtotal: 5000,
				tax: 500,
				taxRate: 10,
				total: 5500
			},
			payment: {
				method: 'Bank Transfer',
				accountName: 'Your Company',
				bsb: '123-456',
				accountNumber: '12345678'
			}
		},
		{
			id: 'INV-251103002',
			status: 'send',
			date: '2025-10-28',
			due: '2025-11-15',
			client: {
				name: 'Tech Startup Inc',
				email: 'billing@techstartup.com'
			},
			provider: {
				name: 'Your Company',
				email: 'you@company.com'
			},
			items: [],
			pricing: {
				subtotal: 3200,
				tax: 320,
				taxRate: 10,
				total: 3520
			},
			payment: {
				method: 'Bank Transfer',
				accountName: 'Your Company',
				bsb: '123-456',
				accountNumber: '12345678'
			}
		},
		{
			id: 'INV-251103003',
			status: 'draft',
			date: '2025-10-15',
			due: '2025-11-01',
			client: {
				name: 'Global Services Ltd',
				email: 'accounts@global.com'
			},
			provider: {
				name: 'Your Company',
				email: 'you@company.com'
			},
			items: [],
			pricing: {
				subtotal: 7500,
				tax: 750,
				taxRate: 10,
				total: 8250
			},
			payment: {
				method: 'Bank Transfer',
				accountName: 'Your Company',
				bsb: '123-456',
				accountNumber: '12345678'
			}
		}
	];

	let searchQuery = '';
	let activeTab: InvoiceStatus | 'all' = 'all';

	$: filteredInvoices = mockInvoices.filter((invoice) => {
		// Filter by status tab
		const matchesStatus = activeTab === 'all' || invoice.status === activeTab;

		// Filter by search query
		const matchesSearch =
			searchQuery === '' ||
			invoice.id.toLowerCase().includes(searchQuery.toLowerCase()) ||
			invoice.client.name.toLowerCase().includes(searchQuery.toLowerCase());

		return matchesStatus && matchesSearch;
	});
</script>

<div class="flex w-full flex-col gap-4">
	<!-- Search bar -->
	<div class="relative">
		<SearchIcon class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
		<Input
			type="text"
			placeholder="Search invoices by number or client..."
			bind:value={searchQuery}
			class="pl-10"
		/>
	</div>

	<!-- Status tabs -->
	<Tabs.Root bind:value={activeTab}>
		<Tabs.List class="grid w-full grid-cols-3">
			<Tabs.Trigger value="all">All</Tabs.Trigger>
			<Tabs.Trigger value="draft">Drafts</Tabs.Trigger>
			<Tabs.Trigger value="send">Sent</Tabs.Trigger>
		</Tabs.List>

		<Tabs.Content value={activeTab} class="mt-4">
			{#if filteredInvoices.length === 0}
				<div class="flex flex-col items-center justify-center rounded-lg border border-dashed p-8">
					<p class="text-muted-foreground">No invoices found</p>
				</div>
			{:else}
				<div class="flex flex-col gap-3">
					{#each filteredInvoices as invoice (invoice.id)}
						<ShelfItem {invoice} />
					{/each}
				</div>
			{/if}
		</Tabs.Content>
	</Tabs.Root>
</div>
