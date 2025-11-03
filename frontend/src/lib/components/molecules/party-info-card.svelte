<!--
	Party Info Card Molecule
	
	Displays provider or client information in a card format.
	Used for displaying contact details on invoices and forms.
	
	Props:
	- party: Party - The party (provider/client) data to display
	- title?: string - Card title (e.g., "Provider", "Bill To")
	- class?: string - Additional CSS classes
	
	Usage:
	<PartyInfoCard party={invoice.provider} title="From" />
	<PartyInfoCard party={invoice.client} title="Bill To" />
-->
<script lang="ts">
	import type { Party } from '@/types/invoice';
	import { cn } from '@/utils';
	import { formatPhone, formatABN } from '@/utils/formatters';

	interface Props {
		party: Party;
		title?: string;
		class?: string;
	}

	let { party, title, class: customClass = '' }: Props = $props();
</script>

<div class={cn('space-y-2', customClass)}>
	{#if title}
		<h3 class="text-sm font-semibold tracking-wide text-muted-foreground uppercase">{title}</h3>
	{/if}

	<div class="space-y-1">
		<p class="font-semibold text-foreground">{party.name}</p>

		{#if party.address}
			<p class="text-sm whitespace-pre-line text-muted-foreground">{party.address}</p>
		{/if}

		{#if party.abn}
			<p class="text-sm text-muted-foreground">
				<span class="font-medium">ABN:</span>
				{formatABN(party.abn)}
			</p>
		{/if}

		{#if party.phone}
			<p class="text-sm text-muted-foreground">
				<span class="font-medium">Phone:</span>
				{formatPhone(party.phone)}
			</p>
		{/if}

		{#if party.email}
			<p class="text-sm text-muted-foreground">
				<span class="font-medium">Email:</span>
				{party.email}
			</p>
		{/if}
	</div>
</div>
