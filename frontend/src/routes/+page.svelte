<script lang="ts">
	import InvoiceShelf from '@/components/organisms/shelf/invoice-shelf.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import { filteredInvoices } from '@/stores';

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
			<InvoiceShelf data={$filteredInvoices} />
		</div>
	</div>
</div>
