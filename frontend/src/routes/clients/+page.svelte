<script lang="ts">
	import ProfileShelf from '@/components/organisms/shelf/profile-shelf.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import { clients } from '@/stores';
	import PlusIcon from '@lucide/svelte/icons/plus';

	function createNewClient() {
		window.location.href = '/clients/new';
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

			<ProfileShelf data={$clients} />
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
