<script lang="ts">
	import InvoiceDisplayCard from '@/components/organisms/invoice-display/invoice-display.svelte';
	// this route is for backend to fetch printable invoice view for PDF generation
	// no edit buttons or other actions here
	// return 500 error if invoice not found

	import { onMount } from 'svelte';
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

	let fontsLoaded = $state(false);

	onMount(async () => {
		if (document.fonts) {
			await document.fonts.ready;
		}
		// Small delay to ensure layout is stable
		setTimeout(() => {
			fontsLoaded = true;
		}, 100);
	});
</script>

{#if error}
	<!-- Error Display -->
	<div id="pdf-render-error" class="container mx-auto max-w-5xl p-4 text-center">
		<p class="text-red-600">{error}</p>
	</div>
{:else if invoice && fontsLoaded}
	<!-- Invoice Display - Only show id="pdf-render-complete" when fonts are ready -->
	<div id="pdf-render-complete" class="container mx-auto max-w-5xl p-4">
		<InvoiceDisplayCard {invoice} class="print:border-none print:shadow-none" />
	</div>
{:else}
	<div class="container mx-auto max-w-5xl p-4 text-center">
		<p>Loading invoice...</p>
	</div>
{/if}
