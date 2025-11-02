<script lang="ts">
	import ProviderEditor from '@/components/custom/profile-editor/provider-editor.svelte';
	import type { ProviderData } from '@/types/invoice';
	import { activeProvider, saveProvider } from '@/stores/provider';
	import { onMount } from 'svelte';

	let currentProvider = $state<ProviderData | null>(null);

	onMount(() => {
		const unsubscribe = activeProvider.subscribe((value) => {
			currentProvider = value;
		});
		return unsubscribe;
	});

	async function handleSave(provider: ProviderData) {
		console.log('Saving provider:', provider);
		await saveProvider(provider);
		window.history.back();
	}

	function handleCancel() {
		window.history.back();
	}
</script>

<div class="container mx-auto max-w-3xl p-4">
	{#if currentProvider}
		<ProviderEditor
			provider={currentProvider}
			mode="edit"
			onSave={handleSave}
			onCancel={handleCancel}
		/>
	{:else}
		<p class="text-center text-muted-foreground">No provider selected</p>
	{/if}
</div>
