<!--
	Invoice Display Card Organism (Refactored from invoice-card.svelte)
	
	Complete invoice display component matching original layout exactly.
	Reduced from 254 lines to ~80 lines by extracting display sections.
	
	Props:
	- invoice: Invoice - Invoice data to display
	- class?: string - Additional CSS classes
	
	Usage:
	<InvoiceDisplayCard {invoice} />
-->
<script lang="ts">
	import type { Invoice } from '@/types/invoice';
	import { cn } from '@/utils';
	import * as Card from '@/components/ui/card';
	import InvoiceDisplayHeader from './invoice-display-header.svelte';
	import InvoiceDisplayParties from './invoice-display-parties.svelte';
	import InvoiceDisplayItemsTable from './invoice-display-items-table.svelte';
	import InvoiceDisplayItemsMobile from './invoice-display-items-mobile.svelte';
	import { PaymentInfoCard, TotalsSummary } from '@/components/molecules';

	interface Props {
		invoice: Invoice;
		class?: string;
	}

	let { invoice, class: customClass = '' }: Props = $props();
</script>

<Card.Root class={cn(customClass)}>
	<Card.Content class="p-4 sm:p-8">
		<!-- Header Section: INVOICE title, invoice number, dates, provider info -->
		<div class="mb-6 sm:mb-8">
			<InvoiceDisplayHeader {invoice} />
		</div>

		<!-- Bill To Section -->
		<div class="mb-6 sm:mb-8">
			<InvoiceDisplayParties client={invoice.client} />
		</div>

		<!-- Invoice Items - Desktop Table -->
		<div class="mb-8">
			<InvoiceDisplayItemsTable items={invoice.items} />
		</div>

		<!-- Invoice Items - Mobile Cards -->
		<div class="mb-6 sm:mb-8">
			<InvoiceDisplayItemsMobile items={invoice.items} />
		</div>

		<!-- Totals Section -->
		<div class="mb-6 flex justify-end sm:mb-8">
			<TotalsSummary pricing={invoice.pricing} class="w-full sm:w-72" />
		</div>

		<!-- Payment Information -->
		<div class="border-t border-border pt-4 sm:pt-6">
			<PaymentInfoCard payment={invoice.payment} />
		</div>

		<!-- Notes/Terms -->
		<div class="mt-6 border-t border-border pt-4 sm:mt-8 sm:pt-6">
			<p class="text-xs text-muted-foreground sm:text-sm">
				Payment is due within 30 days. Please include the invoice number with your payment. Thank
				you for your business!
			</p>
		</div>
	</Card.Content>
</Card.Root>
