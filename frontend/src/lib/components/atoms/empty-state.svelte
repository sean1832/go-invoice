<!--
	Empty State Atom
	
	Simple component for displaying empty states with icon and message.
	
	Props:
	- message: string - The message to display
	- icon?: Snippet - Icon snippet (optional)
	- class?: string - Additional CSS classes
	
	Usage:
	<EmptyState message="No invoices found" />
	<EmptyState message="No items">
		{#snippet icon()}
			<PackageIcon class="mb-4 h-12 w-12 text-muted-foreground opacity-50" />
		{/snippet}
	</EmptyState>
-->
<script lang="ts">
	import { cn } from '@/utils';
	import InboxIcon from '@lucide/svelte/icons/inbox';
	import type { Snippet } from 'svelte';

	interface Props {
		message: string;
		icon?: Snippet;
		class?: string;
	}

	let { message, icon, class: customClass = '' }: Props = $props();
</script>

<div class={cn('flex flex-col items-center justify-center py-12 text-center', customClass)}>
	{#if icon}
		{@render icon()}
	{:else}
		<InboxIcon class="mb-4 h-12 w-12 text-muted-foreground opacity-50" />
	{/if}
	<p class="text-muted-foreground">{message}</p>
</div>
