<script lang="ts">
	import Shelf from '@/components/organisms/shelf/shelf.svelte';
	import ProfileItem from '@/components/organisms/shelf/profile-item.svelte';
	import type { Party } from '@/types/invoice';

	interface Props {
		data: Party[];
		searchPlaceholder?: string;
		emptyMessage?: string;
		onError?: (error: Error) => void;
		onEdit: (item: Party) => void;
		onDelete: (item: Party) => void;
		deletingInvoiceId?: string | null;
	}

	const {
		data,
		searchPlaceholder = 'Search profiles...',
		emptyMessage = 'No profiles found',
		onError = (error: Error) => console.error('Profile shelf error:', error),
		onEdit,
		onDelete,
		deletingProfileId: deletingInvoiceId
	} = $props();
</script>

<Shelf
	{data}
	itemComponent={ProfileItem}
	searchFields={['name', 'email']}
	keyField="id"
	{searchPlaceholder}
	{emptyMessage}
	{onError}
	{onEdit}
	{onDelete}
	deletingItemId={deletingInvoiceId}
/>
