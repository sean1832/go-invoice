<script lang="ts">
	import { ClientForm } from '@/components/organisms/profile-form';
	import { api } from '@/services';
	import type { ClientData } from '@/types/invoice';

	let isSaving = $state(false);
	let saveError = $state<string | null>(null);

	async function handleSave(client: ClientData) {
		isSaving = false;
		saveError = null;

		try {
			await api.clients.createClient(fetch, client);
			window.history.back();
		} catch (err) {
			console.error('failed to save client: ', err);
			saveError = err instanceof Error ? err.message : 'An unknown error occurred while saving.';
		} finally {
			isSaving = false;
		}
	}

	function handleCancel() {
		window.history.back();
	}
</script>

<div class="container mx-auto max-w-3xl p-4">
	<ClientForm mode="create" disable={isSaving} onSave={handleSave} onCancel={handleCancel} />
</div>
