<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from 'mode-watcher';
	import { Toaster } from 'svelte-sonner';
	import { onMount } from 'svelte';
	import { api } from '@/services';
	let { children } = $props();
	import Navbar from '@/components/organisms/navbar/navbar.svelte';

	onMount(async () => {
		// Check authentication status on app load
		try {
			await api.auth.checkSession(fetch);
		} catch (error) {
			console.error('Failed to check auth session:', error);
			// Silently fail - user will see login prompt if needed
		}
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
