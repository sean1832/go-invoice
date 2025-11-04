<script lang="ts">
	import { ProviderForm } from '@/components/organisms/profile-form';
	import type { ProviderData } from '@/types/invoice';
	//import { saveProvider } from '@/stores/provider';

	async function handleSave(provider: ProviderData) {
		console.log('Creating new provider:', provider);
		// TODO: implement save api logic
		//await saveProvider(provider);
		// Navigate back
		window.history.back();
	}

	function handleCancel() {
		window.history.back();
	}
</script>

<div class="container mx-auto max-w-3xl p-4">
	<ProviderForm mode="create" onSave={handleSave} onCancel={handleCancel} />
</div>
