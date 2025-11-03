<!--
	Invoice Line Item Display Molecule
	
	Display-only line item component for showing invoice items.
	Used in invoice display/preview modes.
	
	Props:
	- item: ServiceItem - The service item data to display
	- class?: string - Additional CSS classes
	
	Usage:
	<InvoiceLineItemDisplay item={serviceItem} />
-->
<script lang="ts">
	import type { ServiceItem } from '@/types/invoice';
	import { cn } from '@/utils';
	import CurrencyDisplay from '@/components/atoms/currency-display.svelte';
	import DateDisplay from '@/components/atoms/date-display.svelte';

	interface Props {
		item: ServiceItem;
		class?: string;
	}

	let { item, class: customClass = '' }: Props = $props();
</script>

<div class={cn('border-b py-4 last:border-b-0', customClass)}>
	<div class="grid grid-cols-12 gap-4">
		<!-- Date -->
		<div class="col-span-12 md:col-span-2">
			<DateDisplay date={item.date} format="short" class="text-sm text-muted-foreground" />
		</div>

		<!-- Description -->
		<div class="col-span-12 md:col-span-5">
			<p class="font-semibold text-foreground">{item.description}</p>
			{#if item.descriptionDetail}
				<p class="mt-1 text-sm text-muted-foreground">{item.descriptionDetail}</p>
			{/if}
		</div>

		<!-- Quantity -->
		<div class="col-span-4 text-center md:col-span-1">
			<p class="text-sm text-muted-foreground">{item.quantity}</p>
		</div>

		<!-- Unit Price -->
		<div class="col-span-4 text-right md:col-span-2">
			<CurrencyDisplay amount={item.unitPrice} class="text-sm text-muted-foreground" />
		</div>

		<!-- Total Price -->
		<div class="col-span-4 text-right md:col-span-2">
			<CurrencyDisplay amount={item.totalPrice} class="font-semibold text-foreground" />
		</div>
	</div>
</div>
