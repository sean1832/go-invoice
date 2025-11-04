<script lang="ts">
	import ErrorAlert from '@/components/molecules/error-alert.svelte';
	import ProfileShelf from '@/components/organisms/shelf/profile-shelf.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import { api } from '@/services';
	import { clients } from '@/stores';
	import { removeClient } from '@/stores/clients';
	import type { ClientData } from '@/types/invoice';
	import PlusIcon from '@lucide/svelte/icons/plus';

	function createNewClient() {
		window.location.href = '/clients/new';
	}

	let isLoading = $state(false);
	let errorMessage = $state<string | null>(null);
	let hasLoaded = $state(false);

	async function loadClients() {
		isLoading = true;
		errorMessage = null;

		try {
			const data = await api.clients.getAllClients(fetch);
			console.log(data);
			clients.set(data);
			hasLoaded = true;
		} catch (error) {
			console.error('Failed to load clients:', error);
			errorMessage =
				error instanceof Error ? error.message : 'Failed to load clients. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	// Load clients on mount only
	$effect(() => {
		if (!hasLoaded) {
			loadClients();
		}
	});

	let deletingClientId = $state<string | null>(null);

	async function handleDelete(item: ClientData) {
		deletingClientId = item.id;
		try {
			await api.clients.deleteClient(fetch, item.id);
			// Remove from store - UI updates automatically
			removeClient(item.id);
		} catch (err) {
			console.error('failed to delete client: ', err);
			errorMessage = err instanceof Error ? err.message : 'Unknown error deleting client.';
		} finally {
			deletingClientId = null;
		}
	}

	function handleEdit(item: ClientData) {
		window.location.href = `/clients/${item.id}/edit`;
	}
</script>

<div class="p-4">
	<div class="container mx-auto flex justify-center">
		<div class="w-full max-w-4xl">
			<div class="mb-6 flex items-end justify-between gap-4">
				<div>
					<h1 class="text-3xl font-bold tracking-tight">Clients</h1>
					<p class="text-muted-foreground">Manage and track your clients</p>
				</div>
				<!-- Desktop Create Button -->
				<Button onclick={createNewClient} class="hidden md:flex" size="default">
					<PlusIcon class="mr-2 h-4 w-4" />
					New Client
				</Button>
			</div>

			{#if errorMessage}
				<div class="mb-4">
					<ErrorAlert
						message={errorMessage}
						title="Loading error"
						onRetry={() => {
							errorMessage = null;
						}}
					/>
				</div>
			{:else}
				<ProfileShelf
					data={$clients}
					onEdit={handleEdit}
					onDelete={handleDelete}
					deletingProfileId={deletingClientId}
				/>
			{/if}
		</div>
	</div>

	<!-- Mobile Floating Action Button -->
	<Button
		onclick={createNewClient}
		size="lg"
		class="fixed right-8 bottom-8 h-14 w-14 rounded-full shadow-lg transition-shadow hover:shadow-xl md:hidden"
		title="Create new client"
	>
		<PlusIcon class="h-6 w-6" />
	</Button>
</div>
