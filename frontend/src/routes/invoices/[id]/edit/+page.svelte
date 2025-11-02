<script lang="ts">
	import { page } from '$app/state';
	import InvoiceForm from '@/components/custom/invoice-editor/invoice-editor.svelte';
	import type { Invoice } from '@/types/invoice';

	const invoiceId = page.params.id || 'INV-251103001';

	// Mock data - will be replaced with API call to fetch the invoice
	const mockInvoice: Invoice = {
		id: invoiceId,
		status: 'draft',
		date: '2025-11-01',
		due: '2025-11-30',
		provider: {
			name: 'Your Company',
			address: '123 Business Street',
			email: 'contact@yourcompany.com',
			phone: '+61 3 1234 5678',
			abn: '12 345 678 901'
		},
		client: {
			name: 'Acme Corporation',
			address: '456 Client Avenue',
			email: 'billing@acme.com',
			phone: '+61 2 9876 5432',
			abn: '98 765 432 109'
		},
		items: [
			{
				date: '2025-10-15',
				description: 'Web Development Services',
				quantity: 40,
				unitPrice: 100,
				totalPrice: 4000
			}
		],
		pricing: {
			subtotal: 4000,
			tax: 400,
			taxRate: 10,
			total: 4400
		},
		payment: {
			method: 'Bank Transfer',
			accountName: 'Your Company Pty Ltd',
			bsb: '123-456',
			accountNumber: '12345678'
		}
	};

	let invoice = $state(mockInvoice);

	function handleSave(data: any) {
		// TODO: Implement API call to update invoice
		console.log('Updating invoice:', invoice.id, data);
		// Redirect to invoice view after successful save
		window.location.href = `/invoices/${invoice.id}`;
	}

	function handleCancel() {
		window.location.href = `/invoices/${invoice.id}`;
	}
</script>

<div class="container mx-auto p-4 max-w-5xl">
	<InvoiceForm mode="edit" invoice={invoice} onSave={handleSave} onCancel={handleCancel} />
</div>
