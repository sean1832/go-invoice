<script lang="ts">
	import ProfileShelf from '@/components/custom/shelf/profile-shelf.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import type { ClientData } from '@/types/invoice';
	import PlusIcon from '@lucide/svelte/icons/plus';

	const mockupClients: ClientData[] = [
		{
			id: 'acme_corp',
			name: 'Acme Corporation',
			email: 'contact@acme.com',
			address: '123 Business St',
			phone: '+61 2 9876 5432',
			abn: '98 765 432 109',
			taxRate: 10
		},
		{
			id: 'tech_startup',
			name: 'Tech Startup Inc',
			email: 'info@techstartup.com',
			address: '456 Innovation Way',
			phone: '+61 3 1234 5678',
			abn: '12 345 678 901',
			taxRate: 10
		},
		{
			id: 'global_solutions',
			name: 'Global Solutions Ltd',
			email: 'contact@globalsolutions.com',
			address: '789 Enterprise Rd',
			phone: '+61 7 8765 4321',
			abn: '23 456 789 012',
			taxRate: 10
		}
	];
	function createNewClient() {
		window.location.href = '/clients/new';
	}
</script>

<div class="p-4">
	<div class="container mx-auto flex justify-center">
		<div class="w-full max-w-4xl">
			<div class="mb-6">
				<h1 class="text-3xl font-bold tracking-tight">Clients</h1>
				<p class="text-muted-foreground">Manage and track your clients</p>
			</div>
			<ProfileShelf data={mockupClients} />
		</div>
	</div>

	<!-- Floating Action Button -->
	<Button
		onclick={createNewClient}
		size="lg"
		class="fixed right-8 bottom-8 h-14 w-14 rounded-full shadow-lg transition-shadow hover:shadow-xl"
		title="Create new client"
	>
		<PlusIcon class="h-6 w-6" />
	</Button>
</div>
