<script lang="ts">
	import ProviderSelector from '@/components/custom/provider-selector.svelte';
	import Button from '@/components/ui/button/button.svelte';
</script>

<nav class="p-4">
	<div class="container mx-auto flex justify-between">
		<h1 class="text-2xl font-bold">GoInvoice</h1>
		<div class="flex items-center space-x-4">
			<a href="/"><Button variant="ghost" class="cursor-pointer">Invoices</Button></a>
			<a href="/clients"><Button variant="ghost" class="cursor-pointer">Clients</Button></a>
			<ProviderSelector />
		</div>
	</div>
	<!-- Divider -->
	<div class="my-4 border-t border-gray-800"></div>
</nav>
