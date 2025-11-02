<script lang="ts">
	import InvoiceShelf from '@/components/custom/shelf/invoice-shelf.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import type { Invoice } from '@/types/invoice';

	function createNewInvoice() {
		window.location.href = '/invoices/new';
	}

	// Mock data - this will be replaced with real API calls later
	const mockInvoices: Invoice[] = [
		{
			id: 'INV-251103001',
			status: 'draft',
			date: '2025-11-01',
			due: '2025-11-30',
			client: {
				id: 'acme_corp',
				name: 'Acme Corporation',
				email: 'contact@acme.com',
				address: '123 Business St'
			},
			provider: {
				id: 'your_company',
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
				id: 'tech_startup',
				name: 'Tech Startup Inc',
				email: 'billing@techstartup.com'
			},
			provider: {
				id: 'your_company',
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
				id: 'global_solutions',
				name: 'Global Services Ltd',
				email: 'accounts@global.com'
			},
			provider: {
				id: 'your_company',
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
</script>

<div class="p-4">
	<div class="container mx-auto flex justify-center">
		<div class="w-full max-w-4xl">
			<div class="mb-6">
				<h1 class="text-3xl font-bold tracking-tight">Invoices</h1>
				<p class="text-muted-foreground">Manage and track your invoices</p>
			</div>
			<InvoiceShelf data={mockInvoices} />
		</div>
	</div>

	<!-- Floating Action Button -->
	<Button
		onclick={createNewInvoice}
		size="lg"
		class="fixed right-8 bottom-8 h-14 w-14 rounded-full shadow-lg transition-shadow hover:shadow-xl"
		title="Create new invoice"
	>
		<PlusIcon class="h-6 w-6" />
	</Button>
</div>
