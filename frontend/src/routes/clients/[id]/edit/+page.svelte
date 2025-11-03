<script lang="ts">
	import ClientEditor from '@/components/organisms/profile-form/client-editor.svelte';
	import type { ClientData } from '@/types/invoice';

	// Mock client data - in real app, this would be loaded from the API using the [id] param
	const existingClient: ClientData = {
		id: 'acme_corporation',
		name: 'Acme Corporation',
		email: 'contact@acme.com',
		address: '123 Business St, Sydney NSW 2000',
		phone: '+61 2 9876 5432',
		abn: '98 765 432 109',
		taxRate: 10,
		targetEmail: 'billing@acme.com'
	};

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
	<ClientEditor client={existingClient} mode="edit" onSave={handleSave} onCancel={handleCancel} />
</div>
