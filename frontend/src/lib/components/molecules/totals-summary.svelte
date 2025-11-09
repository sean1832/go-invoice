<!--
	Totals Summary Molecule
	
	Displays invoice pricing breakdown matching original invoice-card.svelte.
	Shows subtotal, tax with borders, and total in colored background.
	
	Props:
	- pricing: Pricing - The pricing object with subtotal, tax, taxRate, total
	- class?: string - Additional CSS classes
	
	Usage:
	<TotalsSummary pricing={invoice.pricing} />
-->
<script lang="ts">
	import type { Pricing } from '@/types/invoice';
	import { cn } from '@/utils';
	import { formatCurrency } from '@/helpers';

	interface Props {
		pricing: Pricing;
		class?: string;
	}

	let { pricing, class: customClass = '' }: Props = $props();
</script>

<div class={cn('w-full', customClass)}>
	<div class="flex justify-between border-b border-border py-2 text-sm sm:text-base">
		<span class="text-muted-foreground">Subtotal:</span>
		<span class="font-semibold text-foreground">{formatCurrency(pricing.subtotal)}</span>
	</div>
	<div class="flex justify-between border-b border-border py-2 text-sm sm:text-base">
		<span class="text-muted-foreground">GST ({pricing.tax_rate}%):</span>
		<span class="font-semibold text-foreground">{formatCurrency(pricing.tax)}</span>
	</div>
	<div
		class="mt-2 flex justify-between rounded-lg bg-primary px-4 py-3 text-base text-primary-foreground sm:text-lg"
	>
		<span class="font-bold">Total:</span>
		<span class="font-bold">{formatCurrency(pricing.total)}</span>
	</div>
</div>
