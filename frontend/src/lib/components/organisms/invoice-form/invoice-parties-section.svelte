<!--
	Invoice Parties Section Organism
	
	Section for provider and client selection in invoice forms.
	Uses ProfileSelector components side-by-side.
	
	Props:
	- provider: Party - Provider data (bindable)
	- client: Party - Client data (bindable)
	- providers: Profile[] - Available providers
	- clients: Profile[] - Available clients
	- onProviderSelect: (id) => void - Callback when provider selected
	- onClientSelect: (id) => void - Callback when client selected
	- class?: string - Additional CSS classes
	
	Usage:
	<InvoicePartiesSection
		bind:provider
		bind:client
		{providers}
		{clients}
		onProviderSelect={handleProviderSelect}
		onClientSelect={handleClientSelect}
	/>
-->
<script lang="ts">
	import type { Party } from '@/types/invoice';
	import { cn } from '@/utils';
	import ProfileSelector from '@/components/custom/invoice-editor/profile-selector.svelte';

	interface Profile {
		id: string;
		name: string;
		email?: string;
		abn?: string;
		address?: string;
		phone?: string;
		taxRate?: number;
		paymentInfo?: any;
	}

	interface Props {
		provider: Party;
		client: Party;
		providers: Profile[];
		clients: Profile[];
		selectedProviderId?: string;
		selectedClientId?: string;
		onProviderSelect: (id: string) => void;
		onClientSelect: (id: string) => void;
		onConfigureProvider?: () => void;
		onAddNewProvider?: () => void;
		onConfigureClient?: () => void;
		onAddNewClient?: () => void;
		class?: string;
	}

	let {
		provider = $bindable(),
		client = $bindable(),
		providers,
		clients,
		selectedProviderId = $bindable(),
		selectedClientId = $bindable(),
		onProviderSelect,
		onClientSelect,
		onConfigureProvider,
		onAddNewProvider,
		onConfigureClient,
		onAddNewClient,
		class: customClass = ''
	}: Props = $props();
</script>

<div class={cn('grid gap-4 md:grid-cols-2', customClass)}>
	<!-- Provider Selector -->
	<ProfileSelector
		class="h-full"
		type="provider"
		profiles={providers}
		bind:selectedProfileId={selectedProviderId}
		profileData={provider}
		onSelect={onProviderSelect}
		onConfigure={onConfigureProvider}
		onAddNew={onAddNewProvider}
	/>

	<!-- Client Selector -->
	<ProfileSelector
		class="h-full"
		type="client"
		profiles={clients}
		bind:selectedProfileId={selectedClientId}
		profileData={client}
		onSelect={onClientSelect}
		onConfigure={onConfigureClient}
		onAddNew={onAddNewClient}
	/>
</div>
