<!--
	Invoice Header Section Organism
	
	Section for invoice ID, issue date, and due date fields.
	Part of the invoice form breakdown.
	
	Props:
	- invoiceId: string - Invoice ID (bindable)
	- issueDate: string - Issue date ISO string (bindable)
	- dueDate: string - Due date ISO string (bindable)
	- class?: string - Additional CSS classes
	
	Usage:
	<InvoiceHeaderSection
		bind:invoiceId
		bind:issueDate
		bind:dueDate
	/>
-->
<script lang="ts">
	import { cn } from '@/utils';
	import { safeParseDate, dateValueToISOString } from '@/helpers';
	import type { DateValue } from '@internationalized/date';
	import * as Card from '@/components/ui/card';
	import Label from '@/components/ui/label/label.svelte';
	import Input from '@/components/ui/input/input.svelte';
	import DatePicker from '@/components/atoms/date-picker.svelte';

	interface Props {
		invoiceId: string;
		issueDate: string;
		dueDate: string;
		class?: string;
	}

	let {
		invoiceId = $bindable(),
		issueDate = $bindable(),
		dueDate = $bindable(),
		class: customClass = ''
	}: Props = $props();

	// Local DateValue states for date pickers
	let issueDateValue = $state<DateValue | undefined>(safeParseDate(issueDate));
	let dueDateValue = $state<DateValue | undefined>(safeParseDate(dueDate));

	// Compute grid columns based on whether invoice ID is present
	const gridCols = $derived(
		invoiceId !== ''
			? 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3' // 3 columns when ID present
			: 'grid-cols-1 sm:grid-cols-2' // 2 columns when ID absent
	);

	// Sync issue date: DateValue <-> string
	$effect(() => {
		if (issueDateValue && dateValueToISOString(issueDateValue) !== issueDate) {
			issueDate = dateValueToISOString(issueDateValue);
		}
	});

	$effect(() => {
		const parsed = safeParseDate(issueDate);
		if (parsed && parsed.toString() !== issueDateValue?.toString()) {
			issueDateValue = parsed;
		}
	});

	// Sync due date: DateValue <-> string
	$effect(() => {
		if (dueDateValue && dateValueToISOString(dueDateValue) !== dueDate) {
			dueDate = dateValueToISOString(dueDateValue);
		}
	});

	$effect(() => {
		const parsed = safeParseDate(dueDate);
		if (parsed && parsed.toString() !== dueDateValue?.toString()) {
			dueDateValue = parsed;
		}
	});
</script>

<Card.Root class={cn(customClass)}>
	<Card.Header>
		<Card.Title>Invoice Details</Card.Title>
		<Card.Description>Basic invoice information and dates</Card.Description>
	</Card.Header>
	<Card.Content>
		<div class={cn('grid gap-4', gridCols)}>
			<!-- Invoice ID -->
			{#if invoiceId !== ''}
				<div class="space-y-2">
					<Label for="invoice-id" class="text-sm font-medium">Invoice ID</Label>
					<Input
						id="invoice-id"
						bind:value={invoiceId}
						placeholder="INV-251103001"
						class="w-full"
					/>
				</div>
			{/if}

			<!-- Issue Date -->
			<div class="space-y-2">
				<Label class="text-sm font-medium">Issue Date</Label>
				<DatePicker bind:value={issueDateValue} placeholder="Select issue date" class="w-full" />
			</div>

			<!-- Due Date -->
			<div class="space-y-2">
				<Label class="text-sm font-medium">Due Date</Label>
				<DatePicker bind:value={dueDateValue} placeholder="Select due date" class="w-full" />
			</div>
		</div>
	</Card.Content>
</Card.Root>
