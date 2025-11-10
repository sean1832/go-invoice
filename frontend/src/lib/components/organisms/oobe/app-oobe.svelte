<!--
	App OOBE (Out-of-Box Experience) Component
	
	Guides new users through initial app setup by checking for required
	provider and client data. Shows a step-by-step checklist with links
	to complete each setup step.
	
	Props:
	- onComplete?: () => void - Callback when all setup steps are completed
-->
<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '@/services';
	import { Button } from '@/components/ui/button';
	import * as Card from '@/components/ui/card';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';
	import CheckCircle2Icon from '@lucide/svelte/icons/check-circle-2';
	import CircleIcon from '@lucide/svelte/icons/circle';
	import ArrowRightIcon from '@lucide/svelte/icons/arrow-right';

	interface Props {
		onComplete?: () => void;
	}

	let { onComplete }: Props = $props();

	let isLoading = $state(true);
	let loadError = $state<string | null>(null);
	let hasProvider = $state(false);
	let hasClient = $state(false);

	async function checkSetupStatus() {
		try {
			isLoading = true;
			loadError = null;

			const [providers, clients] = await Promise.all([
				api.providers.getAllProviders(fetch),
				api.clients.getAllClients(fetch)
			]);

			hasProvider = providers && Array.isArray(providers) && providers.length > 0;
			hasClient = clients && Array.isArray(clients) && clients.length > 0;
		} catch (err) {
			console.error('Failed to check setup status:', err);
			loadError =
				err instanceof Error ? err.message : 'Failed to check setup status. Please try again.';
		} finally {
			isLoading = false;
		}
	}

	onMount(() => {
		checkSetupStatus();
	});

	const steps = $derived([
		{
			id: 'provider',
			title: 'Set up your provider profile',
			description: 'Add your business information that will appear on invoices',
			completed: hasProvider,
			link: '/providers/new'
		},
		{
			id: 'client',
			title: 'Add at least one client',
			description: 'Create a client profile to issue invoices to',
			completed: hasClient,
			link: '/clients/new'
		}
	]);

	const allStepsCompleted = $derived(hasProvider && hasClient);

	function handleContinue() {
		onComplete?.();
	}
</script>

<div class="mx-auto max-w-3xl space-y-6">
	{#if loadError}
		<ErrorAlert message={loadError} showBackButton={false} />
	{:else if isLoading}
		<div class="flex items-center justify-center py-12">
			<p class="text-muted-foreground">Checking setup status...</p>
		</div>
	{:else}
		<!-- Header -->
		<div class="space-y-2">
			<h1 class="text-3xl font-bold tracking-tight">Welcome to Invoice Manager</h1>
			<p class="text-muted-foreground">
				Let's get you set up! Complete these steps to start managing your invoices.
			</p>
		</div>

		<!-- Setup Steps -->
		<div class="space-y-4">
			{#each steps as step}
				<Card.Root class="p-6">
					<div class="flex items-start gap-4">
						<!-- Status Icon -->
						<div class="mt-1">
							{#if step.completed}
								<CheckCircle2Icon class="h-6 w-6 text-green-600" />
							{:else}
								<CircleIcon class="h-6 w-6 text-muted-foreground" />
							{/if}
						</div>

						<!-- Content -->
						<div class="flex-1 space-y-2">
							<h3 class="text-lg font-semibold">{step.title}</h3>
							<p class="text-sm text-muted-foreground">{step.description}</p>

							{#if !step.completed}
								<div class="pt-2">
									<Button href={step.link} variant="outline" size="sm">
										{step.id === 'provider' ? 'Set up provider' : 'Add client'}
										<ArrowRightIcon class="ml-2 h-4 w-4" />
									</Button>
								</div>
							{/if}
						</div>
					</div>
				</Card.Root>
			{/each}
		</div>

		<!-- Continue Button -->
		{#if allStepsCompleted}
			<div class="flex justify-end pt-4">
				<Button onclick={handleContinue} size="lg">
					Continue to dashboard
					<ArrowRightIcon class="ml-2 h-5 w-5" />
				</Button>
			</div>
		{:else}
			<div class="rounded-lg border border-dashed border-muted-foreground/25 bg-muted/50 p-4">
				<p class="text-center text-sm text-muted-foreground">
					Complete the steps above to start using the app
				</p>
			</div>
		{/if}
	{/if}
</div>
