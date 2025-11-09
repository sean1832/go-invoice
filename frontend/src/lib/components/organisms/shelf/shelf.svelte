<script lang="ts">
	import Input from '@/components/ui/input/input.svelte';
	import SearchIcon from '@lucide/svelte/icons/search';

	interface Props {
		data: any[];
		children: any; // Svelte snippet for rendering each item
		searchFields?: string[]; // Fields to search in (e.g., ['name', 'email'])
		keyField?: string; // Field to use as the unique key (default: 'id')
		showTabs?: boolean; // Whether to show status tabs
		statusField?: string; // Field name for status filtering
		tabOptions?: { label: string; value: string }[]; // Tab configuration
		emptyMessage?: string; // Message when no results
		searchPlaceholder?: string;
		deletingItemId?: string | null; // ID of the item currently being deleted
	}

	// Generic props with defaults
	const {
		data,
		children,
		searchFields = [],
		keyField = 'id',
		showTabs = false,
		statusField = undefined,
		tabOptions = [],
		emptyMessage = 'No items found',
		searchPlaceholder = 'Search...',
		deletingItemId = null
	}: Props = $props();

	let searchQuery = $state('');
	let activeTab = $state<string>('all');

	const filteredData = $derived(
		data.filter((item: any) => {
			// Filter by status tab if enabled
			if (showTabs && statusField) {
				const matchesStatus = activeTab === 'all' || item[statusField] === activeTab;
				if (!matchesStatus) return false;
			}

			// Filter by search query
			if (searchQuery === '') return true;

			const query = searchQuery.toLowerCase();
			return searchFields.some((field: string) => {
				const value = item[field];
				return value && String(value).toLowerCase().includes(query);
			});
		})
	);
</script>

<div class="flex w-full flex-col gap-4">
	<!-- Search bar -->
	<div class="relative">
		<SearchIcon class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
		<Input type="text" placeholder={searchPlaceholder} bind:value={searchQuery} class="pl-10" />
	</div>

	<!-- Content -->
	{#if filteredData.length === 0}
		<div class="flex flex-col items-center justify-center rounded-lg border border-dashed p-8">
			<p class="text-muted-foreground">{emptyMessage}</p>
		</div>
	{:else}
		<div class="flex flex-col gap-3">
			{#each filteredData as item, index (item[keyField] || index)}
				{@render children(item, deletingItemId ? deletingItemId === item[keyField] : false)}
			{/each}
		</div>
	{/if}
</div>
