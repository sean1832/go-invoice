<script lang="ts">
	import Button from '@/components/ui/button/button.svelte';
	import Input from '@/components/ui/input/input.svelte';
	import Label from '@/components/ui/label/label.svelte';
	import * as Card from '@/components/ui/card';
	import type { ClientData } from '@/types/invoice';
	import SaveIcon from '@lucide/svelte/icons/save';
	import Spinner from '@/components/atoms/spinner.svelte';
	import XIcon from '@lucide/svelte/icons/x';

	interface Props {
		client?: ClientData;
		onSave?: (client: ClientData) => void;
		onCancel?: () => void;
		disable?: boolean;
		mode?: 'create' | 'edit';
	}

	const { client = undefined, onSave, onCancel, disable = false, mode = 'create' } = $props();

	// Form state
	let formData = $state<ClientData>(
		client || {
			id: '',
			name: '',
			email: '',
			address: '',
			phone: '',
			abn: '',
			tax_rate: 0,
			email_target: '',
			url: ''
		}
	);

	// Validation errors
	let errors = $state<Record<string, string>>({});

	function validateForm(): boolean {
		const newErrors: Record<string, string> = {};

		if (!formData.name?.trim()) {
			newErrors.name = 'Client name is required';
		}

		if (formData.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
			newErrors.email = 'Invalid email format';
		}

		if (formData.tax_rate < 0 || formData.tax_rate > 100) {
			newErrors.taxRate = 'Tax rate must be between 0 and 100';
		}

		if (formData.email_target && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email_target)) {
			newErrors.targetEmail = 'Invalid email format';
		}

		if (formData.abn && formData.abn.length != 11) {
			newErrors.abn = 'ABN must be 11 digits';
		}

		errors = newErrors;
		return Object.keys(newErrors).length === 0;
	}

	function handleSave() {
		if (validateForm()) {
			// Generate ID if creating new client
			if (mode === 'create' && !formData.id) {
				formData.id = formData.name.toLowerCase().replace(/\s+/g, '_');
			}
			onSave?.(formData);
		}
	}

	function handleCancel() {
		onCancel?.();
	}
</script>

<Card.Root class="w-full">
	<Card.Header>
		<Card.Title>{mode === 'create' ? 'New Client' : 'Edit Client'}</Card.Title>
		<Card.Description>
			{mode === 'create'
				? 'Enter client details to create a new profile'
				: 'Update client information'}
		</Card.Description>
	</Card.Header>

	<Card.Content class="space-y-6">
		<!-- Basic Information Section -->
		<div class="space-y-4">
			<h3 class="text-sm font-semibold text-foreground">Basic Information</h3>

			<div class="space-y-2">
				<Label for="name">Client Name *</Label>
				<Input
					id="name"
					type="text"
					placeholder="Acme Corporation"
					bind:value={formData.name}
					disabled={disable}
					class={errors.name ? 'border-destructive' : ''}
				/>
				{#if errors.name}
					<p class="text-sm text-destructive">{errors.name}</p>
				{/if}
			</div>

			<div class="space-y-2">
				<Label for="email">Email</Label>
				<Input
					id="email"
					type="email"
					placeholder="contact@acme.com"
					bind:value={formData.email}
					disabled={disable}
					class={errors.email ? 'border-destructive' : ''}
				/>
				{#if errors.email}
					<p class="text-sm text-destructive">{errors.email}</p>
				{/if}
			</div>

			<div class="grid gap-4 md:grid-cols-2">
				<div class="space-y-2">
					<Label for="phone">Phone</Label>
					<Input
						id="phone"
						type="tel"
						placeholder="+61 2 1234 5678"
						bind:value={formData.phone}
						disabled={disable}
					/>
				</div>

				<div class="space-y-2">
					<Label for="abn">ABN</Label>
					<Input
						id="abn"
						type="text"
						placeholder="12 345 678 901"
						bind:value={formData.abn}
						disabled={disable}
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
					placeholder="123 Business Street, Sydney NSW 2000"
					disabled={disable}
					bind:value={formData.address}
				/>
			</div>

			<div class="space-y-2">
				<Label for="website">Website</Label>
				<Input
					id="website"
					type="text"
					placeholder="https://www.acme.com"
					disabled={disable}
					bind:value={formData.url}
				/>
			</div>
		</div>

		<!-- Invoice Settings Section -->
		<div class="space-y-4">
			<h3 class="text-sm font-semibold text-foreground">Invoice Settings</h3>

			<div class="grid gap-4 md:grid-cols-2">
				<div class="space-y-2">
					<Label for="taxRate">Tax Rate (%) *</Label>
					<Input
						id="taxRate"
						type="number"
						min="0"
						max="100"
						step="0.1"
						placeholder="10"
						bind:value={formData.tax_rate}
						disabled={disable}
						class={errors.taxRate ? 'border-destructive' : ''}
					/>
					{#if errors.taxRate}
						<p class="text-sm text-destructive">{errors.taxRate}</p>
					{:else}
						<p class="text-sm text-muted-foreground">Default tax rate for invoices</p>
					{/if}
				</div>

				<div class="space-y-2">
					<Label for="targetEmail">Invoice Email</Label>
					<Input
						id="targetEmail"
						type="email"
						placeholder="billing@acme.com"
						bind:value={formData.email_target}
						disabled={disable}
						class={errors.targetEmail ? 'border-destructive' : ''}
					/>
					{#if errors.targetEmail}
						<p class="text-sm text-destructive">{errors.targetEmail}</p>
					{:else}
						<p class="text-sm text-muted-foreground">Email for sending invoices (optional)</p>
					{/if}
				</div>
			</div>
		</div>
	</Card.Content>

	<Card.Footer class="flex justify-end gap-3">
		<Button variant="outline" onclick={handleCancel} disabled={disable}>
			<XIcon class="mr-2 h-4 w-4" />
			Cancel
		</Button>
		<Button onclick={handleSave} disabled={disable}>
			{#if disable}
				<Spinner class="mr-2 h-4 w-4" />
			{:else}
				<SaveIcon class="mr-2 h-4 w-4" />
			{/if}
			{mode === 'create' ? 'Create Client' : 'Save Changes'}
		</Button>
	</Card.Footer>
</Card.Root>
