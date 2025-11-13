<script lang="ts">
	import Button from '@/components/ui/button/button.svelte';
	import Input from '@/components/ui/input/input.svelte';
	import Label from '@/components/ui/label/label.svelte';
	import * as Card from '@/components/ui/card';
	import type { ProviderData } from '@/types/invoice';
	import SaveIcon from '@lucide/svelte/icons/save';
	import XIcon from '@lucide/svelte/icons/x';

	interface Props {
		provider?: ProviderData;
		onSave?: (provider: ProviderData) => void;
		onCancel?: () => void;
		mode?: 'create' | 'edit';
	}

	const { provider = undefined, onSave, onCancel, mode = 'create' }: Props = $props();

	// Form state
	let formData = $state<ProviderData>(
		provider || {
			id: '',
			name: '',
			email: '',
			address: '',
			phone: '',
			abn: '',
			payment_info: {
				method: 'EFT Payment',
				account_name: '',
				bsb: '',
				account_number: ''
			}
		}
	);

	// Validation errors
	let errors = $state<Record<string, string>>({});

	function validateForm(): boolean {
		const newErrors: Record<string, string> = {};

		if (!formData.name?.trim()) {
			newErrors.name = 'Provider name is required';
		}

		if (!formData.email?.trim()) {
			newErrors.email = 'Email is required';
		} else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
			newErrors.email = 'Invalid email format';
		}

		if (!formData.payment_info.account_name?.trim()) {
			newErrors.accountName = 'Account name is required';
		}

		const bsb = formData.payment_info.bsb?.trim();
		if (!bsb) {
			newErrors.bsb = 'BSB is required';
		} else if (!/^\d{3}-?\d{3}$/.test(bsb)) {
			newErrors.bsb = 'Invalid BSB format. Must be 6 digits (e.g., 123456 or 123-456)';
		}

		const accountNumber = formData.payment_info.account_number?.trim();
		if (!accountNumber) {
			newErrors.accountNumber = 'Account number is required';
		} else if (!/^\d{6,9}$/.test(accountNumber)) {
			newErrors.accountNumber = 'Invalid account number. Must be 6 to 9 digits';
		}

		const abn = formData.abn?.trim();
		if (abn && abn.length != 11) {
			newErrors.abn = 'ABN must be 11 digits';
		}

		errors = newErrors;
		return Object.keys(newErrors).length === 0;
	}

	function handleSave() {
		if (validateForm()) {
			// trim whitespace for all relevant fields
			formData.payment_info.bsb = formData.payment_info.bsb.replace(/-/g, '').trim();
			formData.payment_info.account_number = formData.payment_info.account_number.trim();
			formData.abn = formData.abn?.trim();
			formData.phone = formData.phone?.trim();
			formData.email = formData.email?.trim();

			onSave?.(formData);
		}
	}

	function handleCancel() {
		onCancel?.();
		// If no cancel callback, go back
		if (!onCancel) {
			window.history.back();
		}
	}
</script>

<Card.Root class="w-full">
	<Card.Header>
		<Card.Title>{mode === 'create' ? 'New Provider' : 'Edit Provider'}</Card.Title>
		<Card.Description>
			{mode === 'create'
				? 'Enter your business details to create a provider profile'
				: 'Update your business information'}
		</Card.Description>
	</Card.Header>

	<Card.Content class="space-y-6">
		<!-- Basic Information Section -->
		<div class="space-y-4">
			<h3 class="text-sm font-semibold text-foreground">Business Information</h3>

			<div class="space-y-2">
				<Label for="name">Business Name *</Label>
				<Input
					id="name"
					type="text"
					placeholder="Your Company Pty Ltd"
					bind:value={formData.name}
					class={errors.name ? 'border-destructive' : ''}
				/>
				{#if errors.name}
					<p class="text-sm text-destructive">{errors.name}</p>
				{/if}
			</div>

			<div class="space-y-2">
				<Label for="email">Email *</Label>
				<Input
					id="email"
					type="email"
					placeholder="hello@yourcompany.com"
					bind:value={formData.email}
					class={errors.email ? 'border-destructive' : ''}
				/>
				{#if errors.email}
					<p class="text-sm text-destructive">{errors.email}</p>
				{/if}
			</div>

			<div class="grid gap-4 md:grid-cols-2">
				<div class="space-y-2">
					<Label for="phone">Phone</Label>
					<Input id="phone" type="tel" placeholder="+61 2 1234 5678" bind:value={formData.phone} />
				</div>

				<div class="space-y-2">
					<Label for="abn">ABN</Label>
					<Input
						id="abn"
						type="text"
						placeholder="12 345 678 901"
						bind:value={formData.abn}
						class={errors.abn ? 'border-destructive' : ''}
					/>
					{#if errors.abn}
						<p class="text-sm text-destructive">{errors.abn}</p>
					{/if}
				</div>
			</div>

			<div class="space-y-2">
				<Label for="address">Address</Label>
				<Input
					id="address"
					type="text"
					placeholder="456 Business Road, Melbourne VIC 3000"
					bind:value={formData.address}
				/>
			</div>
		</div>

		<!-- Payment Information Section -->
		<div class="space-y-4">
			<h3 class="text-sm font-semibold text-foreground">Payment Information</h3>
			<p class="text-sm text-muted-foreground">
				This information will appear on invoices for client payments
			</p>

			<div class="space-y-2">
				<Label for="method">Payment Method</Label>
				<Input
					id="method"
					type="text"
					placeholder="Bank Transfer"
					bind:value={formData.payment_info.method}
				/>
			</div>

			<div class="space-y-2">
				<Label for="accountName">Account Name *</Label>
				<Input
					id="accountName"
					type="text"
					placeholder="Your Company Pty Ltd"
					bind:value={formData.payment_info.account_name}
					class={errors.accountName ? 'border-destructive' : ''}
				/>
				{#if errors.accountName}
					<p class="text-sm text-destructive">{errors.accountName}</p>
				{/if}
			</div>

			<div class="grid gap-4 md:grid-cols-2">
				<div class="space-y-2">
					<Label for="bsb">BSB *</Label>
					<Input
						id="bsb"
						type="text"
						placeholder="123-456"
						bind:value={formData.payment_info.bsb}
						class={errors.bsb ? 'border-destructive' : ''}
					/>
					{#if errors.bsb}
						<p class="text-sm text-destructive">{errors.bsb}</p>
					{/if}
				</div>

				<div class="space-y-2">
					<Label for="accountNumber">Account Number *</Label>
					<Input
						id="accountNumber"
						type="text"
						placeholder="12345678"
						bind:value={formData.payment_info.account_number}
						class={errors.accountNumber ? 'border-destructive' : ''}
					/>
					{#if errors.accountNumber}
						<p class="text-sm text-destructive">{errors.accountNumber}</p>
					{/if}
				</div>
			</div>
		</div>
	</Card.Content>

	<Card.Footer class="flex justify-end gap-3">
		<Button variant="outline" onclick={handleCancel}>
			<XIcon class="mr-2 h-4 w-4" />
			Cancel
		</Button>
		<Button onclick={handleSave}>
			<SaveIcon class="mr-2 h-4 w-4" />
			{mode === 'create' ? 'Create Provider' : 'Save Changes'}
		</Button>
	</Card.Footer>
</Card.Root>
