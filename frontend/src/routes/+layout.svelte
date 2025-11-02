<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from 'mode-watcher';
	import { initializeMockProviders } from '@/stores/provider';
	import { onMount } from 'svelte';
	let { children } = $props();
	import Navbar from '@/components/custom/navbar.svelte';

	onMount(() => {
		initializeMockProviders();
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<ModeWatcher />
<Navbar />
{@render children()}
