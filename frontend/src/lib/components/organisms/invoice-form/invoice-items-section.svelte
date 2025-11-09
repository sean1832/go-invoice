<!--
	Invoice Items Section Organism
	
	Section for managing invoice line items with add/remove functionality.
	Includes line items list and totals summary.
	
	Props:
	- items: ServiceItem[] - Array of service items (bindable)
	- taxRate: number - Tax rate percentage
	- onAddItem: () => void - Callback to add new item
	- onRemoveItem: (index) => void - Callback to remove item
	- onUpdateItem: (index, field, value) => void - Callback to update item
	- class?: string - Additional CSS classes
	
	Usage:
	<InvoiceItemsSection
		bind:items
		{taxRate}
		onAddItem={handleAddItem}
		onRemoveItem={handleRemoveItem}
		onUpdateItem={handleUpdateItem}
	/>
-->
<script lang="ts">
	import type { ServiceItem } from '@/types/invoice';
	import { cn } from '@/utils';
	import { calculatePricing } from '@/helpers';
	import * as Card from '@/components/ui/card';
	import Button from '@/components/ui/button/button.svelte';
	import Separator from '@/components/ui/separator/separator.svelte';
	import InvoiceLineItemEditor from '@/components/molecules/invoice-line-item-editor.svelte';
	import TotalsSummary from '@/components/molecules/totals-summary.svelte';
	import PlusIcon from '@lucide/svelte/icons/plus';

	interface Props {
		items: ServiceItem[];
		taxRate: number;
		onAddItem: () => void;
		onRemoveItem: (index: number) => void;
		onUpdateItem: (index: number, field: keyof ServiceItem, value: any) => void;
		class?: string;
	}

	let {
		items = $bindable(),
		taxRate,
		onAddItem,
		onRemoveItem,
		onUpdateItem,
		class: customClass = ''
	}: Props = $props();

	// Calculate pricing from items
	const pricing = $derived(calculatePricing(items, taxRate));
</script>

<Card.Root class={cn(customClass)}>
	<Card.Header>
		<div class="flex items-center justify-between">
			<div>
				<Card.Title>Service Items</Card.Title>
				<Card.Description>Add items and services for this invoice</Card.Description>
			</div>
			<Button variant="outline" size="sm" onclick={onAddItem}>
				<PlusIcon class="mr-2 h-4 w-4" />
				Add Item
			</Button>
		</div>
	</Card.Header>
	<Card.Content class="space-y-4">
		<!-- Line Items -->
		{#each items as item, index (index)}
			<InvoiceLineItemEditor
				{item}
				{index}
				canRemove={items.length > 1}
				onUpdate={(field, value) => onUpdateItem(index, field, value)}
				onRemove={() => onRemoveItem(index)}
			/>
		{/each}

		<Separator />

		<!-- Totals -->
		<div class="flex justify-end">
			<TotalsSummary {pricing} class="w-64" />
		</div>
	</Card.Content>
</Card.Root>
