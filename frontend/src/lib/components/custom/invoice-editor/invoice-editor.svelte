<script lang="ts">
	import * as Card from '@/components/ui/card';
	import Button from '@/components/ui/button/button.svelte';
	import Input from '@/components/ui/input/input.svelte';
	import Textarea from '@/components/ui/textarea/textarea.svelte';
	import Label from '@/components/ui/label/label.svelte';
	import Separator from '@/components/ui/separator/separator.svelte';
	import type { Invoice, Party, ServiceItem } from '@/types/invoice';
	import PlusIcon from '@lucide/svelte/icons/plus';
	import TrashIcon from '@lucide/svelte/icons/trash-2';
	import SaveIcon from '@lucide/svelte/icons/save';
	import SendIcon from '@lucide/svelte/icons/send';
	import ProfileSelector from './profile-selector.svelte';
	import DatePicker from '../date-picker/date-picker.svelte';
	import { parseDate, type DateValue } from '@internationalized/date';
	import { activeProvider, providers, clients } from '@/stores';
	import { onMount } from 'svelte';

	interface Props {
		invoice?: Invoice;
		mode: 'create' | 'edit';
		onSave?: (data: any) => void;
		onCancel?: () => void;
	}

	let { invoice, mode, onSave, onCancel }: Props = $props();

	// Initialize selected IDs - will be updated when stores load
	let selectedProviderId = $state<string | undefined>(undefined);
	let selectedClientId = $state<string | undefined>(undefined);

	// Form state
	let formData = $state({
		id: invoice?.id || generateInvoiceId(),
		date: invoice?.date || new Date().toISOString().split('T')[0],
		due: invoice?.due || getDefaultDueDate(),
		provider: invoice?.provider || {
			id: '',
			name: '',
			address: '',
			email: '',
			phone: '',
			abn: ''
		},
		client: invoice?.client || {
			id: '',
			name: '',
			address: '',
			email: '',
			phone: '',
			abn: ''
		},
		items: invoice?.items || [createEmptyLineItem()],
		taxRate: invoice?.pricing?.taxRate || 10
	});

	// Payment info from selected provider (read-only)
	let paymentInfo = $state({
		method: 'Bank Transfer',
		accountName: '',
		bsb: '',
		accountNumber: ''
	});

	// Initialize from activeProvider when it's available
	$effect(() => {
		if ($activeProvider && !invoice && mode === 'create') {
			// Set provider data from active provider
			if (!formData.provider.id) {
				formData.provider = {
					id: $activeProvider.id,
					name: $activeProvider.name,
					email: $activeProvider.email || '',
					abn: $activeProvider.abn || '',
					address: $activeProvider.address || '',
					phone: $activeProvider.phone || ''
				};
			}
			// Set selected provider ID
			if (!selectedProviderId) {
				selectedProviderId = $activeProvider.id;
			}
			// Set payment info
			if ($activeProvider.paymentInfo) {
				paymentInfo = { ...$activeProvider.paymentInfo };
			}
		}
	});

	// DateValue states for date pickers
	let issueDateValue = $state<DateValue | undefined>(undefined);
	let dueDateValue = $state<DateValue | undefined>(undefined);
	let lineItemDateValues = $state<(DateValue | undefined)[]>([]);

	// Helper to safely parse date strings
	function safeParseDate(str: string | undefined | null): DateValue | undefined {
		if (!str) return undefined;
		try {
			return parseDate(str);
		} catch (e) {
			console.warn(`Invalid date string: "${str}"`, e);
			return undefined;
		}
	}

	// Sync issue date: DateValue -> string
	$effect(() => {
		if (issueDateValue && issueDateValue.toString() !== formData.date) {
			formData.date = issueDateValue.toString();
		}
	});

	// Sync issue date: string -> DateValue
	$effect(() => {
		if (formData.date && formData.date !== issueDateValue?.toString()) {
			issueDateValue = safeParseDate(formData.date);
		} else if (!formData.date && issueDateValue) {
			issueDateValue = undefined;
		}
	});

	// Sync due date: DateValue -> string
	$effect(() => {
		if (dueDateValue && dueDateValue.toString() !== formData.due) {
			formData.due = dueDateValue.toString();
		}
	});

	// Sync due date: string -> DateValue
	$effect(() => {
		if (formData.due && formData.due !== dueDateValue?.toString()) {
			dueDateValue = safeParseDate(formData.due);
		} else if (!formData.due && dueDateValue) {
			dueDateValue = undefined;
		}
	});

	// Initialize and sync line item dates
	$effect(() => {
		// Ensure lineItemDateValues array matches items length
		if (lineItemDateValues.length !== formData.items.length) {
			lineItemDateValues = formData.items.map((item) => safeParseDate(item.date));
		}
	});

	// Sync line item dates: DateValue -> string
	$effect(() => {
		lineItemDateValues.forEach((dateValue, index) => {
			if (
				dateValue &&
				formData.items[index] &&
				dateValue.toString() !== formData.items[index].date
			) {
				updateLineItem(index, 'date', dateValue.toString());
			}
		});
	});

	// Sync line item dates: string -> DateValue
	$effect(() => {
		formData.items.forEach((item, index) => {
			if (item.date && item.date !== lineItemDateValues[index]?.toString()) {
				lineItemDateValues[index] = safeParseDate(item.date);
			}
		});
	});

	// Computed values
	let subtotal = $derived(formData.items.reduce((sum, item) => sum + item.totalPrice, 0));
	let tax = $derived(subtotal * (formData.taxRate / 100));
	let total = $derived(subtotal + tax);

	// Helper functions
	function generateInvoiceId(): string {
		const date = new Date();
		const year = date.getFullYear().toString().slice(-2);
		const month = (date.getMonth() + 1).toString().padStart(2, '0');
		const day = date.getDate().toString().padStart(2, '0');
		const random = Math.floor(Math.random() * 1000)
			.toString()
			.padStart(3, '0');
		return `INV-${year}${month}${day}${random}`;
	}

	function getDefaultDueDate(): string {
		const date = new Date();
		date.setDate(date.getDate() + 30);
		return date.toISOString().split('T')[0];
	}

	function createEmptyLineItem(): ServiceItem {
		return {
			date: new Date().toISOString().split('T')[0],
			description: '',
			quantity: 1,
			unitPrice: 0,
			totalPrice: 0
		};
	}

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-AU', {
			style: 'currency',
			currency: 'AUD'
		}).format(amount);
	}

	// Profile selection handlers
	function handleProviderSelect(providerId: string) {
		const provider = $providers.find((p) => p.id === providerId);
		if (provider) {
			selectedProviderId = providerId;
			formData.provider = {
				id: provider.id,
				name: provider.name,
				email: provider.email || '',
				abn: provider.abn || '',
				address: provider.address || '',
				phone: provider.phone || ''
			};
			// Update payment info from provider
			if (provider.paymentInfo) {
				paymentInfo = { ...provider.paymentInfo };
			}
		}
	}

	function handleClientSelect(clientId: string) {
		const client = $clients.find((c) => c.id === clientId);
		if (client) {
			selectedClientId = clientId;
			formData.client = {
				id: client.id,
				name: client.name,
				email: client.email || '',
				abn: client.abn || '',
				address: client.address || '',
				phone: client.phone || ''
			};
			// Update tax rate from client
			if (client.taxRate !== undefined) {
				formData.taxRate = client.taxRate;
			}
		}
	}

	function handleConfigureProvider() {
		console.log('Navigate to provider settings');
		window.location.href = '/settings';
	}

	function handleAddNewProvider() {
		console.log('Navigate to add new provider');
		window.location.href = '/providers/new';
	}

	function handleConfigureClient() {
		console.log('Navigate to client settings');
		window.location.href = `/clients/${selectedClientId}/edit`;
	}

	function handleAddNewClient() {
		console.log('Navigate to add new client');
		window.location.href = '/clients/new';
	}

	// Line item handlers
	function addLineItem() {
		formData.items = [...formData.items, createEmptyLineItem()];
	}

	function removeLineItem(index: number) {
		if (formData.items.length > 1) {
			formData.items = formData.items.filter((_, i) => i !== index);
		}
	}

	function updateLineItem(index: number, field: keyof ServiceItem, value: any) {
		const item = { ...formData.items[index] };
		(item as any)[field] = value;

		// Recalculate totalPrice
		if (field === 'quantity' || field === 'unitPrice') {
			item.totalPrice = item.quantity * item.unitPrice;
		}

		formData.items[index] = item;
		formData.items = [...formData.items];
	}

	// Form actions
	function handleSaveDraft() {
		const invoiceData: Invoice = {
			id: formData.id,
			status: 'draft',
			date: formData.date,
			due: formData.due,
			provider: formData.provider,
			client: formData.client,
			items: formData.items,
			pricing: {
				subtotal,
				tax,
				taxRate: formData.taxRate,
				total
			},
			payment: paymentInfo
		};
		console.log('Save as draft:', invoiceData);
		onSave?.(invoiceData);
	}

	function handleSend() {
		const invoiceData: Invoice = {
			id: formData.id,
			status: 'send',
			date: formData.date,
			due: formData.due,
			provider: formData.provider,
			client: formData.client,
			items: formData.items,
			pricing: {
				subtotal,
				tax,
				taxRate: formData.taxRate,
				total
			},
			payment: paymentInfo
		};
		console.log('Send invoice:', invoiceData);
		onSave?.(invoiceData);
	}

	function handleCancel() {
		// if (confirm('Are you sure you want to cancel? Unsaved changes will be lost.')) {
		// 	onCancel?.();
		// }
		window.location.href = '/';
	}
</script>

<div class="space-y-6">
	<!-- Invoice Header -->
	<Card.Root>
		<Card.Header>
			<Card.Title>
				{mode === 'create' ? 'Create New Invoice' : 'Edit Invoice'}
			</Card.Title>
			<Card.Description>
				{mode === 'create'
					? 'Fill in the details to create a new invoice'
					: 'Update the invoice information'}
			</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-6">
			<!-- Basic Info -->
			<div class="space-y-2">
				<Label>Invoice Number</Label>
				<Input bind:value={formData.id} placeholder="INV-25110301" />
			</div>

			<!-- Dates -->
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label>Issue Date</Label>
					<DatePicker bind:value={issueDateValue} placeholder="select an invoice date" />
				</div>

				<div class="space-y-2">
					<Label>Due Date</Label>
					<DatePicker bind:value={dueDateValue} placeholder="select a due date" />
				</div>
			</div>
		</Card.Content>
	</Card.Root>

	<!-- Provider Information (Read-only) -->
	<ProfileSelector
		class="h-full"
		type="provider"
		profiles={$providers}
		bind:selectedProfileId={selectedProviderId}
		profileData={formData.provider}
		onSelect={handleProviderSelect}
		onConfigure={handleConfigureProvider}
		onAddNew={handleAddNewProvider}
	/>

	<!-- Client Information -->
	<ProfileSelector
		class="h-full"
		type="client"
		profiles={$clients}
		bind:selectedProfileId={selectedClientId}
		profileData={formData.client}
		onSelect={handleClientSelect}
		onConfigure={handleConfigureClient}
		onAddNew={handleAddNewClient}
	/>

	<!-- Line Items -->
	<Card.Root>
		<Card.Header>
			<div class="flex items-center justify-between">
				<div>
					<Card.Title>Service Items</Card.Title>
					<Card.Description>Add items and services for this invoice</Card.Description>
				</div>
				<Button variant="outline" size="sm" onclick={addLineItem}>
					<PlusIcon class="mr-2 h-4 w-4" />
					Add Item
				</Button>
			</div>
		</Card.Header>
		<Card.Content class="space-y-4">
			{#each formData.items as item, index}
				<div class="space-y-4 rounded-lg border p-4">
					<div class="flex items-start gap-2">
						<div class="flex-1 space-y-4">
							<div class="space-y-2">
								<Label>Date</Label>
								<DatePicker bind:value={lineItemDateValues[index]} placeholder="select a date" />
							</div>

							<div class="space-y-2">
								<Label>Description</Label>
								<Textarea
									value={item.description}
									oninput={(e) => updateLineItem(index, 'description', e.currentTarget.value)}
									placeholder="Service or product description"
									rows={2}
								/>
							</div>

							<div class="grid grid-cols-3 gap-4">
								<div class="space-y-2">
									<Label>Quantity</Label>
									<Input
										type="number"
										value={item.quantity}
										oninput={(e) =>
											updateLineItem(index, 'quantity', parseFloat(e.currentTarget.value) || 0)}
										min="0"
										step="1"
									/>
								</div>

								<div class="space-y-2">
									<Label>Unit Price</Label>
									<Input
										type="number"
										value={item.unitPrice}
										oninput={(e) =>
											updateLineItem(index, 'unitPrice', parseFloat(e.currentTarget.value) || 0)}
										min="0"
										step="0.01"
									/>
								</div>

								<div class="space-y-2">
									<Label>Amount</Label>
									<Input
										value={formatCurrency(item.totalPrice)}
										disabled
										class="bg-muted font-semibold"
									/>
								</div>
							</div>
						</div>

						<Button
							variant="ghost"
							size="sm"
							onclick={() => removeLineItem(index)}
							disabled={formData.items.length === 1}
							class="mt-8"
						>
							<TrashIcon class="h-4 w-4 text-destructive" />
						</Button>
					</div>
				</div>
			{/each}

			<Separator />

			<!-- Totals -->
			<div class="flex justify-end">
				<div class="w-64 space-y-2">
					<div class="flex justify-between text-sm">
						<span class="text-muted-foreground">Subtotal:</span>
						<span>{formatCurrency(subtotal)}</span>
					</div>
					<div class="flex justify-between text-sm">
						<span class="text-muted-foreground">Tax ({formData.taxRate}%):</span>
						<span>{formatCurrency(tax)}</span>
					</div>
					<Separator />
					<div class="flex justify-between text-lg font-bold">
						<span>Total:</span>
						<span>{formatCurrency(total)}</span>
					</div>
				</div>
			</div>
		</Card.Content>
	</Card.Root>

	<!-- Action Buttons -->
	<div class="flex justify-between">
		<Button variant="outline" onclick={handleCancel}>Cancel</Button>

		<div class="flex gap-2">
			<Button variant="outline" onclick={handleSaveDraft}>
				<SaveIcon class="mr-2 h-4 w-4" />
				Save as Draft
			</Button>
			<Button variant="default" onclick={handleSend}>
				<SendIcon class="mr-2 h-4 w-4" />
				Save & Send
			</Button>
		</div>
	</div>
</div>
