<script lang="ts">
	import LoaderCircleIcon from '@lucide/svelte/icons/loader-circle';

	interface Props {
		size?: number;
		class?: string;
	}

	let { size = 24, class: className = '' }: Props = $props();
</script>

<div class="flex items-center justify-center {className}">
	<LoaderCircleIcon class="animate-spin" {size} />
</div>
