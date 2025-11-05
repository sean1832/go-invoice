<script lang="ts">
	import InvoiceDisplayCard from '@/components/organisms/invoice-display/invoice-display-card.svelte';
	// this route is for backend to fetch printable invoice view for PDF generation
	// no edit buttons or other actions here
	// return 500 error if invoice not found

	import type { Invoice } from '@/types/invoice';

	interface Props {
		data: {
			invoice: Invoice | null;
			error?: string;
		};
	}
	let { data }: Props = $props();
	let invoice = $derived(data.invoice as Invoice);
	let error = $derived(data.error);

	$effect(() => {
		if (error) {
			throw new Error(error);
		}
	});
</script>

{#if invoice}
	<!-- Printable Invoice Display Component -->
	<div class="container mx-auto max-w-5xl p-4">
		<InvoiceDisplayCard {invoice} />
	</div>
{:else}
	<div class="container mx-auto max-w-5xl p-4 text-center">
		<p class="text-red-600">Invoice not found.</p>
	</div>
{/if}
