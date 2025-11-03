<script lang="ts">
	import { page } from '$app/state';
	import InvoiceForm from '@/components/organisms/invoice-form/invoice-form.svelte';
	import { mockInvoices } from '@/stores';

	const invoiceId = page.params.id || 'INV-251103001';

	let invoice = $state(mockInvoices.find((inv) => inv.id === invoiceId) || mockInvoices[0]);

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

<div class="container mx-auto max-w-5xl p-4">
	<InvoiceForm mode="edit" {invoice} onSave={handleSave} onCancel={handleCancel} />
</div>
