<script lang="ts">
	import { page } from '$app/state';
	import InvoiceForm from '@/components/organisms/invoice-form/invoice-form.svelte';
	import { mockInvoices } from '@/stores';
	import type { Invoice } from '@/types/invoice';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';

	interface Props {
		data: {
			invoice: Invoice | null;
			error?: string;
		};
	}

	let { data }: Props = $props();
	let invoice = $derived(data.invoice as Invoice);
	let error = $derived(data.error);

	function handleSave(data: any) {
		// TODO: Implement API call to update invoice
		console.log('Updating invoice:', invoice.id, data);
		// Redirect to invoice view after successful save
		window.location.href = `/`;
	}

	function handleCancel() {
		window.history.back();
	}
</script>

<div class="container mx-auto max-w-5xl p-4">
	{#if error || !invoice}
		<ErrorAlert message={error || 'Invoice not found'} showBackButton={true} />
	{:else}
		<InvoiceForm mode="edit" {invoice} onSave={handleSave} onCancel={handleCancel} />
	{/if}
</div>
