<!--
	Invoice Display Parties Organism
	
	"Bill To" section matching original invoice-card.svelte layout.
	Shows client information in a muted background card.
	
	Props:
	- client: Party - Client information
	- class?: string - Additional CSS classes
	
	Usage:
	<InvoiceDisplayParties client={invoice.client} />
-->
<script lang="ts">
	import type { Party } from '@/types/invoice';
	import { cn } from '@/utils';

	interface Props {
		client: Party;
		class?: string;
	}

	let { client, class: customClass = '' }: Props = $props();
</script>

<div class={cn(customClass)}>
	<h3 class="mb-3 text-base font-semibold text-foreground sm:text-lg">Bill To:</h3>
	<div class="rounded-lg bg-muted p-3 sm:p-4">
		<p class="font-semibold text-foreground">{client.name}</p>
		<div class="mt-1 space-y-0.5 text-sm text-muted-foreground sm:text-base">
			{#if client.address}
				<p class="whitespace-pre-line">{client.address}</p>
			{/if}
			{#if client.abn}
				<p>ABN: {client.abn}</p>
			{/if}
			{#if client.phone}
				<p>Phone: {client.phone}</p>
			{/if}
			{#if client.email}
				<p>Email: {client.email}</p>
			{/if}
		</div>
	</div>
</div>
