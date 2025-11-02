<script lang="ts">
	import Shelf from '@/components/custom/shelf/shelf.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import PlusIcon from '@lucide/svelte/icons/plus';

	function createNewInvoice() {
		window.location.href = '/invoices/new';
	}
</script>

<div class="p-4">
	<div class="container mx-auto flex justify-center">
		<div class="w-full max-w-4xl">
			<div class="mb-6">
				<h1 class="text-3xl font-bold tracking-tight">Invoices</h1>
				<p class="text-muted-foreground">Manage and track your invoices</p>
			</div>
			<Shelf />
		</div>
	</div>

	<!-- Floating Action Button -->
	<Button
		onclick={createNewInvoice}
		size="lg"
		class="fixed bottom-8 right-8 h-14 w-14 rounded-full shadow-lg hover:shadow-xl transition-shadow"
		title="Create new invoice"
	>
		<PlusIcon class="h-6 w-6" />
	</Button>
</div>


