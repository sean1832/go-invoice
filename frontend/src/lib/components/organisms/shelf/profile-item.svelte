<script lang="ts">
	import * as Item from '$lib/components/ui/item/index.js';
	import Button from '@/components/ui/button/button.svelte';
	import type { Party } from '@/types/invoice';
	import EditIcon from '@lucide/svelte/icons/pencil';
	import TrashIcon from '@lucide/svelte/icons/trash-2';

	interface Props {
		item: Party;
		onEdit: (item: Party) => void;
		onDelete: (item: Party) => void;
		isDeleting?: boolean;
	}

	const { item, onEdit, onDelete, isDeleting } = $props();

	function handleEdit() {
		onEdit(item);
	}

	function handleDelete() {
		if (confirm(`Are you sure you want to delete ${item.name}?`)) {
			onDelete(item);
		}
	}
</script>

<Item.Root variant="outline" class="cursor-pointer transition-colors hover:bg-muted/50">
	<Item.Content class="flex-1" onclick={handleEdit}>
		<div class="flex flex-col gap-1">
			<Item.Title class="text-2xl font-semibold">{item.name}</Item.Title>
			{#if item.email}
				<Item.Description>{item.email}</Item.Description>
			{/if}
			{#if item.address}
				<Item.Description class="text-xs">{item.address}</Item.Description>
			{/if}
			{#if item.phone}
				<Item.Description class="text-xs">{item.phone}</Item.Description>
			{/if}
			{#if item.abn}
				<Item.Description class="text-xs">ABN: {item.abn}</Item.Description>
			{/if}
		</div>
	</Item.Content>

	<Item.Actions class="flex gap-1">
		<Button
			size="sm"
			variant="ghost"
			onclick={handleEdit}
			title="Edit profile"
			aria-label="Edit {item.name}"
			disabled={isDeleting}
		>
			<EditIcon class="h-4 w-4" />
		</Button>
		<Button
			size="sm"
			variant="ghost"
			onclick={handleDelete}
			title="Delete profile"
			aria-label="Delete {item.name}"
			disabled={isDeleting}
		>
			<TrashIcon class="h-4 w-4 text-red-500" />
		</Button>
	</Item.Actions>
</Item.Root>
