<script lang="ts">
	import { page } from '$app/state';
	import InvoiceForm from '@/components/organisms/invoice-form/invoice-form.svelte';
	import { mockInvoices } from '@/stores';
	import type { Invoice } from '@/types/invoice';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import Button from '@/components/ui/button/button.svelte';
	import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';

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

	function goBack() {
		window.history.back();
	}
</script>

<div class="container mx-auto max-w-5xl p-4">
	{#if error || !invoice}
		<!-- Error State -->
		<div class="mb-6">
			<Button variant="ghost" size="sm" onclick={goBack} class="mb-4">
				<ArrowLeftIcon class="mr-2 h-4 w-4" />
				Back
			</Button>
			<Alert.Root variant="destructive">
				<AlertCircleIcon class="h-4 w-4" />
				<Alert.Title>Error</Alert.Title>
				<Alert.Description>
					{error || 'Invoice not found'}
				</Alert.Description>
			</Alert.Root>
		</div>
	{:else}
		<InvoiceForm mode="edit" {invoice} onSave={handleSave} onCancel={handleCancel} />
	{/if}
</div>
