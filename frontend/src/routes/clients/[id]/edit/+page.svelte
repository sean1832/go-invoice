<script lang="ts">
	import { ClientForm } from '@/components/organisms/profile-form';
	import type { ClientData } from '@/types/invoice';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';

	interface Props {
		data: {
			client: ClientData | null;
			error?: string;
		};
	}

	let { data }: Props = $props();
	let client = $derived(data.client);
	let error = $derived(data.error);

	function handleSave(client: ClientData) {
		console.log('Saving client:', client);
		// In real app: Save to API, then navigate back
		window.location.href = '/clients';
	}

	function handleCancel() {
		window.location.href = '/clients';
	}
</script>

<div class="container mx-auto max-w-3xl p-4">
	{#if error || !client}
		<ErrorAlert message={error || 'Client not found'} showBackButton={true} />
	{:else}
		<ClientForm {client} mode="edit" onSave={handleSave} onCancel={handleCancel} />
	{/if}
</div>
