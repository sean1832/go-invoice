<script lang="ts">
	import * as Card from '@/components/ui/card';
	import * as Select from '@/components/ui/select';
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

	interface Props {
		invoice?: Invoice;
		mode: 'create' | 'edit';
		onSave?: (data: any) => void;
		onCancel?: () => void;
	}

	let { invoice, mode, onSave, onCancel }: Props = $props();

	// Form state
	let formData = $state({
		id: invoice?.id || generateInvoiceId(),
		date: invoice?.date || new Date().toISOString().split('T')[0],
		due: invoice?.due || getDefaultDueDate(),
		provider: invoice?.provider || {
			name: '',
			address: '',
			email: '',
			phone: '',
			abn: ''
		},
		client: invoice?.client || {
			name: '',
			address: '',
			email: '',
			phone: '',
			abn: ''
		},
		items: invoice?.items || [createEmptyLineItem()],
		taxRate: invoice?.pricing?.taxRate || 10,
		paymentMethod: invoice?.payment?.method || 'Bank Transfer',
		accountName: invoice?.payment?.accountName || '',
		bsb: invoice?.payment?.bsb || '',
		accountNumber: invoice?.payment?.accountNumber || ''
	});

	// Mock data for providers and clients (these would come from API/storage later)
	let providers = $state([
		{ id: '1', name: 'Zeke Zhang', email: 'zeke@example.com', abn: '12 345 678 901' },
		{ id: '2', name: 'Lan Zhang', email: 'lan@example.com', abn: '98 765 432 109' }
	]);

	let clients = $state([
		{ id: '1', name: 'Dingyu Xu', email: 'dingyu@example.com', abn: '11 222 333 444' },
		{ id: '2', name: 'Client Corp', email: 'contact@clientcorp.com', abn: '55 666 777 888' }
	]);

	// Selected IDs for dropdowns
	let selectedProviderId = $state<string | undefined>(undefined);
	let selectedClientId = $state<string | undefined>(undefined);

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
		const provider = providers.find((p) => p.id === providerId);
		if (provider) {
			selectedProviderId = providerId;
			formData.provider = {
				name: provider.name,
				email: provider.email,
				abn: provider.abn,
				address: '', // These would come from full profile
				phone: ''
			};
		}
	}

	function handleClientSelect(clientId: string) {
		const client = clients.find((c) => c.id === clientId);
		if (client) {
			selectedClientId = clientId;
			formData.client = {
				name: client.name,
				email: client.email,
				abn: client.abn,
				address: '', // These would come from full profile
				phone: ''
			};
		}
	}

	function handleConfigureProvider() {
		console.log('Navigate to provider settings');
		// TODO: Navigate to settings page
	}

	function handleAddNewProvider() {
		console.log('Navigate to add new provider');
		// TODO: Navigate to add provider page
	}

	function handleConfigureClient() {
		console.log('Navigate to client settings');
		// TODO: Navigate to settings page
	}

	function handleAddNewClient() {
		console.log('Navigate to add new client');
		// TODO: Navigate to add client page
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
			payment: {
				method: formData.paymentMethod,
				accountName: formData.accountName,
				bsb: formData.bsb,
				accountNumber: formData.accountNumber
			}
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
			payment: {
				method: formData.paymentMethod,
				accountName: formData.accountName,
				bsb: formData.bsb,
				accountNumber: formData.accountNumber
			}
		};
		console.log('Send invoice:', invoiceData);
		onSave?.(invoiceData);
	}

	function handleCancel() {
		if (confirm('Are you sure you want to cancel? Unsaved changes will be lost.')) {
			onCancel?.();
		}
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
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label>Invoice Number</Label>
					<Input bind:value={formData.id} placeholder="INV-25110301" />
				</div>

				<div class="space-y-2">
					<Label>Tax Rate (%)</Label>
					<Input
						type="number"
						bind:value={formData.taxRate}
						placeholder="10"
						min="0"
						max="100"
						step="0.01"
					/>
				</div>
			</div>

			<!-- Dates -->
			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label>Issue Date</Label>
					<Input type="date" bind:value={formData.date} />
				</div>

				<div class="space-y-2">
					<Label>Due Date</Label>
					<Input type="date" bind:value={formData.due} />
				</div>
			</div>
		</Card.Content>
	</Card.Root>

	<!-- Provider Information (Read-only) -->
	<div class="flex gap-4">
		<ProfileSelector
			class="flex-1"
			type="provider"
			profiles={providers}
			bind:selectedProfileId={selectedProviderId}
			profileData={formData.provider}
			onSelect={handleProviderSelect}
			onConfigure={handleConfigureProvider}
			onAddNew={handleAddNewProvider}
		/>

		<!-- Client Information -->
		<ProfileSelector
			class="flex-1"
			type="client"
			profiles={clients}
			bind:selectedProfileId={selectedClientId}
			profileData={formData.client}
			onSelect={handleClientSelect}
			onConfigure={handleConfigureClient}
			onAddNew={handleAddNewClient}
		/>
	</div>

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
								<Input
									type="date"
									value={item.date}
									oninput={(e) => updateLineItem(index, 'date', e.currentTarget.value)}
								/>
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

	<!-- Payment Information -->
	<Card.Root>
		<Card.Header>
			<Card.Title>Payment Information</Card.Title>
			<Card.Description>Bank account details for payment</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<div class="space-y-2">
				<Label>Payment Method</Label>
				<Input bind:value={formData.paymentMethod} placeholder="Bank Transfer" />
			</div>

			<div class="space-y-2">
				<Label>Account Name</Label>
				<Input bind:value={formData.accountName} placeholder="Your Company Pty Ltd" />
			</div>

			<div class="grid grid-cols-2 gap-4">
				<div class="space-y-2">
					<Label>BSB</Label>
					<Input bind:value={formData.bsb} placeholder="123-456" />
				</div>

				<div class="space-y-2">
					<Label>Account Number</Label>
					<Input bind:value={formData.accountNumber} placeholder="12345678" />
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
