<script lang="ts">
	import Shelf from '@/components/custom/shelf/shelf.svelte';
	import ProfileItem from '@/components/custom/shelf/profile-item.svelte';
	import type { Party } from '@/types/invoice';

	interface Props {
		data: Party[];
		searchPlaceholder?: string;
		emptyMessage?: string;
	}

	const {
		data,
		searchPlaceholder = 'Search profiles...',
		emptyMessage = 'No profiles found'
	} = $props();
</script>

<Shelf
	{data}
	itemComponent={ProfileItem}
	searchFields={['name', 'email']}
	keyField="id"
	{searchPlaceholder}
	{emptyMessage}
/>
