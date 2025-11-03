<script lang="ts">
	import { ClientForm } from '@/components/organisms/profile-form';
	import type { ClientData } from '@/types/invoice';
	import { page } from '$app/state';
	import { mockClients } from '@/stores';

	const clientId = page.params.id;
	if (!clientId) {
		throw new Error('Client ID is required for editing');
	}
	// In real app: Fetch client data from API
	const client = mockClients.find((c) => c.id === clientId) || mockClients[0];

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
	<ClientForm client={client} mode="edit" onSave={handleSave} onCancel={handleCancel} />
</div>
