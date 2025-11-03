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
	import { safeParseDate, dateValueToISOString } from '@/utils/date-helpers';
	import type { DateValue } from '@internationalized/date';
	import * as Card from '@/components/ui/card';
	import Label from '@/components/ui/label/label.svelte';
	import Input from '@/components/ui/input/input.svelte';
	import DatePicker from '@/components/custom/date-picker/date-picker.svelte';

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
	<Card.Content class="space-y-4">
		<div class="grid gap-4 md:grid-cols-3">
			<!-- Invoice ID -->
			<div class="space-y-2">
				<Label for="invoice-id">Invoice ID</Label>
				<Input id="invoice-id" bind:value={invoiceId} placeholder="INV-251103001" />
			</div>

			<!-- Issue Date -->
			<div class="space-y-2">
				<Label>Issue Date</Label>
				<DatePicker bind:value={issueDateValue} placeholder="Select issue date" />
			</div>

			<!-- Due Date -->
			<div class="space-y-2">
				<Label>Due Date</Label>
				<DatePicker bind:value={dueDateValue} placeholder="Select due date" />
			</div>
		</div>
	</Card.Content>
</Card.Root>
