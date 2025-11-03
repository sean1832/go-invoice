<!--
	Invoice Display Header Organism
	
	Header section matching original invoice-card.svelte layout.
	Shows "INVOICE" title, invoice number, dates on left, provider info on right.
	
	Props:
	- invoice: Invoice - Invoice data to display
	- class?: string - Additional CSS classes
	
	Usage:
	<InvoiceDisplayHeader {invoice} />
-->
<script lang="ts">
	import type { Invoice } from '@/types/invoice';
	import { cn } from '@/utils';
	import { formatDateShort } from '@/utils/formatters';

	interface Props {
		invoice: Invoice;
		class?: string;
	}

	let { invoice, class: customClass = '' }: Props = $props();
</script>

<div class={cn('flex flex-col gap-6 sm:flex-row sm:items-start sm:justify-between', customClass)}>
	<!-- Left: Invoice Info -->
	<div>
		<h1 class="mb-2 text-3xl font-bold text-foreground sm:text-4xl">INVOICE</h1>
		<div class="space-y-1 text-sm sm:text-base">
			<p class="text-muted-foreground">
				Invoice Number: <span class="font-semibold text-foreground">{invoice.id}</span>
			</p>
			<p class="text-muted-foreground">
				Date: <span class="font-semibold text-foreground">{formatDateShort(invoice.date)}</span>
			</p>
			<p class="text-muted-foreground">
				Due Date: <span class="font-semibold text-foreground">{formatDateShort(invoice.due)}</span>
			</p>
		</div>
	</div>

	<!-- Right: Provider Info -->
	<div class="sm:text-right">
		<h2 class="mb-2 text-xl font-bold text-foreground sm:text-2xl">{invoice.provider.name}</h2>
		<div class="space-y-0.5 text-sm text-muted-foreground sm:text-base">
			{#if invoice.provider.address}
				<p class="whitespace-pre-line">{invoice.provider.address}</p>
			{/if}
			{#if invoice.provider.abn}
				<p>ABN: {invoice.provider.abn}</p>
			{/if}
			{#if invoice.provider.phone}
				<p>Phone: {invoice.provider.phone}</p>
			{/if}
			{#if invoice.provider.email}
				<p>Email: {invoice.provider.email}</p>
			{/if}
		</div>
	</div>
</div>
