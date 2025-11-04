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
	import { formatCurrency } from '@/utils/formatters';
	import { safeParseDate, dateValueToISOString } from '@/utils/date-helpers';
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

<div class={cn('space-y-4 rounded-lg border p-4', customClass)}>
	<div class="flex items-start gap-2">
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
			<div class="grid grid-cols-3 gap-4">
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

				<div class="space-y-2">
					<Label>Amount</Label>
					<Input value={formatCurrency(item.total_price)} disabled class="bg-muted font-semibold" />
				</div>
			</div>
		</div>

		<!-- Remove Button -->
		<Button
			variant="ghost"
			size="sm"
			onclick={onRemove}
			disabled={!canRemove}
			class="mt-8"
			title={canRemove ? 'Remove item' : 'At least one item required'}
		>
			<TrashIcon class="h-4 w-4 text-destructive" />
		</Button>
	</div>
</div>
