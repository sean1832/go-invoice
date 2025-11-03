<!--
	Invoice Display Items Mobile Organism
	
	Mobile card view matching original invoice-card.svelte layout.
	Shows items as individual cards with date, description, quantity/price grid, and total.
	
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
	import Separator from '@/components/ui/separator/separator.svelte';

	interface Props {
		items: ServiceItem[];
		class?: string;
	}

	let { items, class: customClass = '' }: Props = $props();
</script>

<div class={cn('space-y-4 lg:hidden', customClass)}>
	{#each items as item}
		<div
			class="rounded-lg border-2 border-border bg-card p-4 transition-colors hover:border-primary/50"
		>
			<!-- Item Header -->
			<div class="mb-2 flex items-start justify-between">
				<div class="flex-1">
					<p class="mb-1 font-semibold text-foreground">{item.description}</p>
					{#if item.description_detail}
						<p class="mb-2 text-sm text-muted-foreground">{item.description_detail}</p>
					{/if}
					<p class="text-xs text-muted-foreground">{formatDateShort(item.date)}</p>
				</div>
			</div>

			<Separator class="my-3" />

			<!-- Item Details Grid -->
			<div class="grid grid-cols-2 gap-3 text-sm">
				<div>
					<p class="mb-0.5 text-xs text-muted-foreground">Quantity</p>
					<p class="font-semibold text-foreground">{item.quantity}</p>
				</div>
				<div>
					<p class="mb-0.5 text-xs text-muted-foreground">Unit Price</p>
					<p class="font-semibold text-foreground">{formatCurrency(item.unit_price)}</p>
				</div>
			</div>

			<!-- Item Total -->
			<div class="mt-3 border-t border-border pt-3">
				<div class="flex items-center justify-between">
					<p class="font-medium text-muted-foreground">Amount</p>
					<p class="text-lg font-bold text-foreground">{formatCurrency(item.total_price)}</p>
				</div>
			</div>
		</div>
	{/each}
</div>
