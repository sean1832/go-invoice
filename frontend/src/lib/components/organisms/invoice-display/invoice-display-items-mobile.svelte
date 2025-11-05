<!--
	Invoice Display Items Mobile Organism
	
	Improved mobile card view with better visual hierarchy and cleaner design.
	Shows items as cards with description prominent, metadata compact, and amount highlighted.
	
	Props:
	- items: ServiceItem[] - Array of service items to display
	- class?: string - Additional CSS classes
	
	Usage:
	<InvoiceDisplayItemsMobile items={invoice.items} />
-->
<script lang="ts">
	import type { ServiceItem } from '@/types/invoice';
	import { cn } from '@/utils';
	import { formatCurrency, formatDateShort } from '@/utils/formatters';
	import DateDisplay from '@/components/atoms/date-display.svelte';
	import CurrencyDisplay from '@/components/atoms/currency-display.svelte';

	interface Props {
		items: ServiceItem[];
		class?: string;
	}

	let { items, class: customClass = '' }: Props = $props();
</script>

<div class={cn('space-y-3 md:hidden', customClass)}>
	{#each items as item, index}
		<div
			class="rounded-lg border border-border bg-card p-4 shadow-sm transition-all hover:shadow-md"
		>
			<!-- Description & Date Header -->
			<div class="mb-3">
				<div class="mb-1 flex items-start justify-between gap-2">
					<h4 class="flex-1 text-base leading-tight font-semibold text-foreground">
						{item.description}
					</h4>
					<DateDisplay
						date={item.date}
						format="short"
						class="shrink-0 text-xs text-muted-foreground"
					/>
				</div>
				{#if item.description_detail}
					<p class="mt-1.5 text-sm leading-relaxed text-muted-foreground">
						{item.description_detail}
					</p>
				{/if}
			</div>

			<!-- Quantity & Unit Price Row -->
			<div class="mb-3 flex items-center justify-between text-sm">
				<div class="flex items-center gap-4">
					<div>
						<span class="text-muted-foreground">Qty:</span>
						<span class="ml-1 font-medium text-foreground">{item.quantity}</span>
					</div>
					<div class="h-4 w-px bg-border"></div>
					<div>
						<span class="text-muted-foreground">@ </span>
						<CurrencyDisplay amount={item.unit_price} class="font-medium text-foreground" />
					</div>
				</div>
			</div>

			<!-- Amount (highlighted) -->
			<div class="flex items-center justify-between rounded-md bg-primary/5 px-3 py-2.5">
				<span class="text-sm font-medium text-muted-foreground">Amount</span>
				<CurrencyDisplay amount={item.total_price} class="text-lg font-bold text-primary" />
			</div>
		</div>
	{/each}
</div>
