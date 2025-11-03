<script lang="ts">
	import ProviderEditor from '@/components/organisms/profile-form/provider-editor.svelte';
	import type { ProviderData } from '@/types/invoice';
	import { saveProvider } from '@/stores/provider';

	async function handleSave(provider: ProviderData) {
		console.log('Creating new provider:', provider);
		await saveProvider(provider);
		// Navigate back
		window.history.back();
	}

	function handleCancel() {
		window.history.back();
	}
</script>

<div class="container mx-auto max-w-3xl p-4">
	<ProviderEditor mode="create" onSave={handleSave} onCancel={handleCancel} />
</div>
