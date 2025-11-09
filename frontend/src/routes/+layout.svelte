<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from 'mode-watcher';
	import { Toaster } from 'svelte-sonner';
	import { onMount } from 'svelte';
	let { children } = $props();
	import Navbar from '@/components/organisms/navbar/navbar.svelte';

	onMount(async () => {
		// Initialize all stores when the app starts
		//await Promise.all([initializeProviders()]);
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Go-Invoice</title>
</svelte:head>

<ModeWatcher />
<Toaster position="bottom-center" />
<Navbar class="print:hidden" />
{@render children()}
