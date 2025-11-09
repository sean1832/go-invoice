<!--
	Invoice Form Organism (Refactored from invoice-editor.svelte)
	
	Main invoice form component - orchestrates all invoice form sections.
	Significantly reduced from 539 lines to ~150 lines by extracting sections.
	
	Props:
	- invoice?: Invoice - Existing invoice for edit mode
	- mode: 'create' | 'edit' - Form mode
	- onSave?: (invoice) => void - Callback when form saved
	- onCancel?: () => void - Callback when form cancelled
	
	Usage:
	<InvoiceForm mode="create" onSave={handleSave} onCancel={handleCancel} />
	<InvoiceForm {invoice} mode="edit" onSave={handleUpdate} />
-->
<script lang="ts">
	import type { Invoice, Party, ServiceItem } from '@/types/invoice';
	import { providers, clients, activeProvider } from '@/stores';
	import {
		validateInvoice,
		calculateLineItemTotal,
		calculatePricing,
		getDefaultIssueDate,
		getDefaultDueDate,
		createEmptyLineItem,
		createEmptyParty
	} from '@/helpers';
	import InvoiceHeaderSection from './invoice-header-section.svelte';
	import InvoicePartiesSection from './invoice-parties-section.svelte';
	import InvoiceItemsSection from './invoice-items-section.svelte';
	import InvoiceFormActions from './invoice-form-actions.svelte';

	interface Props {
		invoice?: Invoice;
		mode: 'create' | 'edit';
		isSaving?: boolean;
		onSave?: (data: any) => void;
		onCancel?: () => void;
	}

	let { invoice, mode, isSaving = false, onSave, onCancel }: Props = $props();

	// Form state
	let invoiceId = $state(invoice?.id || '');
	let issueDate = $state(invoice?.date || getDefaultIssueDate());
	let dueDate = $state(invoice?.due || getDefaultDueDate());
	let provider = $state<Party>(invoice?.provider || createEmptyParty());
	let client = $state<Party>(invoice?.client || createEmptyParty());
	let items = $state<ServiceItem[]>(invoice?.items || [createEmptyLineItem()]);
	let taxRate = $state(invoice?.pricing?.tax_rate || 10);

	// Selected IDs for profile selectors
	let selectedProviderId = $state<string | undefined>(undefined);
	let selectedClientId = $state<string | undefined>(undefined);

	// Payment info from provider
	let paymentInfo = $state({
		method: 'Bank Transfer',
		account_name: '',
		bsb: '',
		account_number: ''
	});

	// Initialize from activeProvider in create mode
	$effect(() => {
		if ($activeProvider && !invoice && mode === 'create' && !provider.id) {
			provider = {
				id: $activeProvider.id,
				name: $activeProvider.name,
				email: $activeProvider.email || '',
				abn: $activeProvider.abn || '',
				address: $activeProvider.address || '',
				phone: $activeProvider.phone || ''
			};
			selectedProviderId = $activeProvider.id;
			if ($activeProvider.payment_info) {
				paymentInfo = { ...$activeProvider.payment_info };
			}
		} else if (mode === 'edit' && invoice) {
			selectedProviderId = invoice.provider.id;
			selectedClientId = invoice.client.id;
			if (invoice.payment) {
				paymentInfo = { ...invoice.payment };
			}
		}
	});

	// Profile selection handlers
	function handleProviderSelect(providerId: string) {
		const selected = $providers.find((p) => p.id === providerId);
		if (selected) {
			selectedProviderId = providerId;
			provider = {
				id: selected.id,
				name: selected.name,
				email: selected.email || '',
				abn: selected.abn || '',
				address: selected.address || '',
				phone: selected.phone || ''
			};
			if (selected.payment_info) {
				paymentInfo = { ...selected.payment_info };
			}
		}
	}

	function handleClientSelect(clientId: string) {
		const selected = $clients.find((c) => c.id === clientId);
		if (selected) {
			selectedClientId = clientId;
			client = {
				id: selected.id,
				name: selected.name,
				email: selected.email || '',
				abn: selected.abn || '',
				address: selected.address || '',
				phone: selected.phone || ''
			};
			if (selected.tax_rate !== undefined) {
				taxRate = selected.tax_rate;
			}
		}
	}

	// Navigation handlers
	function handleConfigureProvider() {
		window.location.href = '/settings';
	}

	function handleAddNewProvider() {
		window.location.href = '/providers/new';
	}

	function handleConfigureClient() {
		window.location.href = `/clients/${selectedClientId}/edit`;
	}

	function handleAddNewClient() {
		window.location.href = '/clients/new';
	}

	// Line item handlers
	function handleAddItem() {
		items = [...items, createEmptyLineItem()];
	}

	function handleRemoveItem(index: number) {
		if (items.length > 1) {
			items = items.filter((_, i) => i !== index);
		}
	}

	function handleUpdateItem(index: number, field: keyof ServiceItem, value: any) {
		const item = { ...items[index] };
		(item as any)[field] = value;

		// Recalculate totalPrice if quantity or unitPrice changed
		if (field === 'quantity' || field === 'unit_price') {
			item.total_price = calculateLineItemTotal(item.quantity, item.unit_price);
		}

		items[index] = item;
	}

	// Form submission handlers
	function handleSaveDraft() {
		submitForm();
	}

	function handleCancel() {
		onCancel?.();
	}

	function submitForm() {
		// Calculate final pricing
		const pricing = calculatePricing(items, taxRate);

		// Build invoice object
		const invoiceData: Invoice = {
			id: invoiceId,
			date: issueDate,
			due: dueDate,
			provider,
			client,
			items,
			pricing,
			payment: paymentInfo,
			status: 'draft' // Always draft on save
		};

		// Validate
		const validation = validateInvoice(invoiceData);
		if (!validation.isValid) {
			alert('Please fix the following errors:\n\n' + validation.errors.join('\n'));
			return;
		}

		// Submit
		onSave?.(invoiceData);
	}
</script>

<div class="space-y-6">
	<!-- Invoice Header (ID, Dates) -->
	<InvoiceHeaderSection bind:invoiceId bind:issueDate bind:dueDate />

	<!-- Provider & Client Selection -->
	<InvoicePartiesSection
		bind:provider
		bind:client
		providers={$providers}
		clients={$clients}
		bind:selectedProviderId
		bind:selectedClientId
		onProviderSelect={handleProviderSelect}
		onClientSelect={handleClientSelect}
		onConfigureProvider={handleConfigureProvider}
		onAddNewProvider={handleAddNewProvider}
		onConfigureClient={handleConfigureClient}
		onAddNewClient={handleAddNewClient}
	/>

	<!-- Line Items -->
	<InvoiceItemsSection
		bind:items
		{taxRate}
		onAddItem={handleAddItem}
		onRemoveItem={handleRemoveItem}
		onUpdateItem={handleUpdateItem}
	/>

	<!-- Action Buttons -->
	<InvoiceFormActions {isSaving} onCancel={handleCancel} onSaveDraft={handleSaveDraft} />
</div>
