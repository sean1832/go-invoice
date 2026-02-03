<!--
	Invoice Line Item Editor Molecule
	
	Editable line item component for invoice forms.
	Handles date, description, quantity, unit price display and editing.
	
	Props:
	- item: ServiceItem - The service item data
	- index: number - Item index for labeling
	- canRemove?: boolean - Whether remove button should be enabled
	- onUpdate: (field, value) => void - Callback when item field changes
	- onRemove: () => void - Callback when remove button clicked
	- class?: string - Additional CSS classes
	
	Usage:
	<InvoiceLineItemEditor
		{item}
		{index}
		canRemove={items.length > 1}
		onUpdate={(field, value) => updateItem(index, field, value)}
		onRemove={() => removeItem(index)}
	/>
-->
<script lang="ts">
	import type { ServiceItem } from '@/types/invoice';
	import { cn } from '@/utils';
	import { formatCurrency } from '@/helpers';
	import { safeParseDate, dateValueToISOString } from '@/helpers';
	import type { DateValue } from '@internationalized/date';
	import Label from '@/components/ui/label/label.svelte';
	import Input from '@/components/ui/input/input.svelte';
	import Textarea from '@/components/ui/textarea/textarea.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import DatePicker from '@/components/atoms/date-picker.svelte';
	import TrashIcon from '@lucide/svelte/icons/trash-2';

	interface Props {
		item: ServiceItem;
		index: number;
		canRemove?: boolean;
		onUpdate: (field: keyof ServiceItem, value: any) => void;
		onRemove: () => void;
		class?: string;
	}

	let {
		item,
		index,
		canRemove = true,
		onUpdate,
		onRemove,
		class: customClass = ''
	}: Props = $props();

	// Collapsed state for mobile (defaults to false/collapsed)
	let isExpanded = $state(false);

	// Local date state for date picker
	let dateValue = $state<DateValue | undefined>(safeParseDate(item.date));

	// Sync date picker changes back to item
	$effect(() => {
		if (dateValue && dateValue.toString() !== item.date) {
			onUpdate('date', dateValueToISOString(dateValue));
		}
	});

	// Update dateValue when item.date changes externally
	$effect(() => {
		const parsed = safeParseDate(item.date);
		if (parsed && parsed.toString() !== dateValue?.toString()) {
			dateValue = parsed;
		}
	});
</script>

<div class={cn('space-y-4 rounded-lg border p-3 md:p-4', customClass)}>
	<!-- Mobile Collapsed Header -->
	<div
		class="flex cursor-pointer items-center justify-between md:hidden"
		onclick={() => (isExpanded = !isExpanded)}
		role="button"
		tabindex="0"
		onkeydown={(e) => {
			if (e.key === 'Enter' || e.key === ' ') {
				e.preventDefault();
				isExpanded = !isExpanded;
			}
		}}
	>
		<div class="flex flex-col">
			<span class="font-medium">{item.description || 'New Item'}</span>
			<div class="text-sm text-muted-foreground">
				{item.quantity} x {formatCurrency(item.unit_price)}
			</div>
		</div>
		<div class="flex items-center gap-2">
			<span class="font-bold">{formatCurrency(item.total_price)}</span>
			
			<Button 
				variant="ghost" 
				size="sm" 
				class="h-8 w-8 p-0 text-destructive hover:text-destructive"
				onclick={(e) => {
					e.stopPropagation();
					if (confirm('Remove this item?')) {
						onRemove();
					}
				}}
				disabled={!canRemove}
				title="Remove item"
			>
				<TrashIcon class="h-4 w-4" />
			</Button>

			<Button variant="ghost" size="sm" class="h-8 w-8 p-0">
				{#if isExpanded}
					<span class="text-xs">▼</span>
				{:else}
					<span class="text-xs">▶</span>
				{/if}
			</Button>
		</div>
	</div>

	<div
		class={cn(
			'flex flex-col items-start gap-4 md:flex-row md:gap-2',
			!isExpanded && 'hidden md:flex'
		)}
	>
		<div class="flex-1 space-y-4">
			<!-- Date -->
			<div class="space-y-2">
				<Label>Date</Label>
				<DatePicker bind:value={dateValue} placeholder="Select date" />
			</div>

			<!-- Description -->
			<div class="space-y-2">
				<Label>Description</Label>
				<Textarea
					value={item.description}
					oninput={(e) => onUpdate('description', e.currentTarget.value)}
					placeholder="Service or product description"
					rows={2}
				/>
			</div>

			<!-- Quantity, Unit Price, Amount -->
			<div class="grid grid-cols-2 gap-4 md:grid-cols-3">
				<div class="space-y-2">
					<Label>Quantity</Label>
					<Input
						type="number"
						value={item.quantity}
						oninput={(e) => onUpdate('quantity', parseFloat(e.currentTarget.value) || 0)}
						min="0"
						step="1"
					/>
				</div>

				<div class="space-y-2">
					<Label>Unit Price</Label>
					<Input
						type="number"
						value={item.unit_price}
						oninput={(e) => onUpdate('unit_price', parseFloat(e.currentTarget.value) || 0)}
						min="0"
						step="0.01"
					/>
				</div>

				<div class="col-span-2 space-y-2 md:col-span-1">
					<Label>Amount</Label>
					<Input value={formatCurrency(item.total_price)} disabled class="bg-muted font-semibold" />
				</div>
			</div>
		</div>

		<!-- Remove Button -->
		<Button
			variant="ghost"
			size="sm"
			onclick={() => {
				if (confirm('Remove this item?')) {
					onRemove();
				}
			}}
			disabled={!canRemove}
			class="hidden w-full md:mt-8 md:flex md:w-auto"
			title={canRemove ? 'Remove item' : 'At least one item required'}
		>
			<TrashIcon class="h-4 w-4 text-destructive" />
		</Button>
	</div>
</div>
