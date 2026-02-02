<!--
	Invoice Display Items Table Organism
	
	Desktop table view for invoice line items.
	Shows date, description, quantity, unit price, and amount in columns.
	
	Props:
	- items: ServiceItem[] - Array of service items to display
	- class?: string - Additional CSS classes
	
	Usage:
	<InvoiceDisplayItemsTable items={invoice.items} />
-->
<script lang="ts">
	import type { ServiceItem } from '@/types/invoice';
	import { cn } from '@/utils';
	import DateDisplay from '@/components/atoms/date-display.svelte';
	import CurrencyDisplay from '@/components/atoms/currency-display.svelte';

	interface Props {
		items: ServiceItem[];
		class?: string;
	}

	let { items, class: customClass = '' }: Props = $props();
</script>

<div class={cn('hidden overflow-hidden rounded-lg border-2 border-border md:block', customClass)}>
	<table class="w-full">
		<thead>
			<tr class="bg-primary text-primary-foreground">
				<th class="w-36 px-4 py-4 text-left font-semibold">Date</th>
				<th class="px-6 py-4 text-left font-semibold">Description</th>
				<th class="w-28 px-4 py-4 text-center font-semibold">Qty</th>
				<th class="w-36 px-4 py-4 text-right font-semibold">Unit Price</th>
				<th class="w-36 px-4 py-4 text-right font-semibold">Amount</th>
			</tr>
		</thead>
		<tbody class="bg-card">
			{#each items as item, index}
				<tr
					class="transition-colors hover:bg-muted/50 {index < items.length - 1
						? 'border-b-2 border-border' // add border except for last row
						: ''}"
				>
					<td class="px-4 py-4 text-left">
						<DateDisplay date={item.date} format="short" class="text-muted-foreground" />
					</td>
					<td class="px-6 py-4">
						<p class="font-semibold text-foreground">{item.description}</p>
						{#if item.description_detail}
							<p class="mt-1 text-sm text-muted-foreground">{item.description_detail}</p>
						{/if}
					</td>
					<td class="px-6 py-4 text-center text-muted-foreground">
						{item.quantity}
					</td>
					<td class="px-6 py-4 text-right">
						<CurrencyDisplay amount={item.unit_price} class="text-muted-foreground" />
					</td>
					<td class="px-6 py-4 text-right">
						<CurrencyDisplay amount={item.total_price} class="font-semibold text-foreground" />
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
