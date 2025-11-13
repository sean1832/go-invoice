<script lang="ts">
	import ProviderSelector from '@/components/molecules/provider-selector.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import { cn } from '@/utils';

	interface Props {
		class?: string;
	}

	let { class: className = '' }: Props = $props();
</script>

<nav>
	<div class={cn('mt-4 mr-4 ml-4', className)}>
		<div class="container mx-auto flex justify-between">
			<a href="/"
				><h1 class="text-2xl font-bold">
					<img src="/favicon.svg" alt="InvoiceGo Logo" class="mr-2 inline h-10 w-10 align-middle" />
					Go-Invoice
				</h1>
			</a>

			<div class="flex items-center space-x-4">
				<a href="/"
					><Button variant="ghost" class="hidden cursor-pointer md:block">Invoices</Button></a
				>
				<a href="/clients"><Button variant="ghost" class="cursor-pointer">Clients</Button></a>
				<ProviderSelector />
			</div>
		</div>
	</div>
	<!-- Divider -->
	<div class="darkborder-gray-800 my-4 border-t"></div>
</nav>
