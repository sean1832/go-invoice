<script lang="ts">
	import { page } from '$app/state';
	import { Badge } from '@/components/ui/badge';
	import Button from '@/components/ui/button/button.svelte';
	import InvoiceDisplayCard from '@/components/organisms/invoice-display/invoice-display-card.svelte';
	import type { Invoice } from '@/types/invoice';
	import EditIcon from '@lucide/svelte/icons/pencil';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import CopyIcon from '@lucide/svelte/icons/copy';
	import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';
	import { mockInvoices } from '@/stores';

	// Mock invoice data - will be replaced with API call
	// const mockInvoice: Invoice = {
	// 	id: page.params.id || 'INV-251103001',
	// 	status: 'draft',
	// 	date: '2025-11-01',
	// 	due: '2025-12-01',
	// 	provider: {
	// 		id: 'your_company',
	// 		name: 'Your Company',
	// 		address: '123 Business Street, Suite 100\nMelbourne, VIC 3000',
	// 		email: 'contact@yourcompany.com',
	// 		phone: '+61 3 1234 5678',
	// 		abn: '12 345 678 901'
	// 	},
	// 	client: {
	// 		id: 'acme_corporation',
	// 		name: 'Acme Corporation',
	// 		address: '456 Client Avenue\nSydney, NSW 2000',
	// 		email: 'billing@acme.com',
	// 		phone: '+61 2 9876 5432',
	// 		abn: '98 765 432 109'
	// 	},
	// 	items: [
	// 		{
	// 			date: '2025-10-15',
	// 			description: 'Web Development Services',
	// 			descriptionDetail: 'Full-stack development for e-commerce platform',
	// 			quantity: 40,
	// 			unitPrice: 150,
	// 			totalPrice: 6000
	// 		},
	// 		{
	// 			date: '2025-10-22',
	// 			description: 'UI/UX Design',
	// 			descriptionDetail: 'Design mockups and prototypes',
	// 			quantity: 20,
	// 			unitPrice: 120,
	// 			totalPrice: 2400
	// 		},
	// 		{
	// 			date: '2025-10-28',
	// 			description: 'Consulting Services',
	// 			quantity: 10,
	// 			unitPrice: 180,
	// 			totalPrice: 1800
	// 		}
	// 	],
	// 	pricing: {
	// 		subtotal: 10200,
	// 		tax: 1020,
	// 		taxRate: 10,
	// 		total: 11220
	// 	},
	// 	payment: {
	// 		method: 'Bank Transfer',
	// 		accountName: 'Your Company Pty Ltd',
	// 		bsb: '123-456',
	// 		accountNumber: '12345678'
	// 	}
	// };
	const id = page.params.id;
	const inv = mockInvoices.find((inv) => inv.id === id) || mockInvoices[0];

	let invoice = $state(inv);

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
</div>
