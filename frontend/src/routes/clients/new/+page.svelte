<script lang="ts">
	import ClientEditor from '@/components/custom/profile-editor/client-editor.svelte';
	import type { ClientData } from '@/types/invoice';

	function handleSave(client: ClientData) {
		console.log('Creating new client:', client);
		// In real app: Save to API, then navigate to client list or detail page
		window.location.href = '/clients';
	}

	function handleCancel() {
		window.location.href = '/clients';
	}
</script>

<div class="container mx-auto max-w-3xl p-4">
	<ClientEditor mode="create" onSave={handleSave} onCancel={handleCancel} />
</div>
