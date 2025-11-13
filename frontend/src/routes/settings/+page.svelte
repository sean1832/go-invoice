<script lang="ts">
	import Button from '@/components/ui/button/button.svelte';
	import * as Card from '@/components/ui/card';
	import type { ProviderData } from '@/types/invoice';
	import {
		activeProvider,
		authStore,
		isAuthenticated,
		currentUserEmail,
		requiresOAuth
	} from '@/stores';
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import CheckCircle2Icon from '@lucide/svelte/icons/check-circle-2';
	import { ProviderForm } from '@/components/organisms/profile-form';
	import { api } from '@/services';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';
	import { updateProvider } from '@/stores/provider';
	import { onMount } from 'svelte';

	// Use reactive reference to activeProvider
	let currentProvider = $derived($activeProvider);

	let errorMessage = $state<string | null>(null);
	let isSaving = $state(false);
	let isLoggingIn = $state(false);
	let isLoggingOut = $state(false);
	let version = $state<string>('dev');

	onMount(async () => {
		try {
			version = await api.version.getVersion(fetch);
		} catch (err) {
			console.error('Failed to fetch version:', err);
		}
	});

	async function handleSave(provider: ProviderData) {
		isSaving = true;
		try {
			await api.providers.updateProvider(fetch, provider.id, provider);

			// Update both the providers list AND the activeProvider
			updateProvider(provider);

			// Go back after saving
			window.history.back();
		} catch (err) {
			console.error('failed to update provider data: ', err);
			errorMessage = err instanceof Error ? err.message : 'Unknown error updating provider data.';
		} finally {
			isSaving = false;
		}
	}

	function handleCancel() {
		// Go back without saving
		window.history.back();
	}

	function navigateToNewProvider() {
		window.location.href = '/providers/new';
	}

	async function handleLogin() {
		isLoggingIn = true;
		try {
			await api.auth.loginWithGoogle();
			// Session will be updated by loginWithGoogle
		} catch (err) {
			console.error('Login failed:', err);
			errorMessage = err instanceof Error ? err.message : 'Failed to connect Google account';
		} finally {
			isLoggingIn = false;
		}
	}

	async function handleLogout() {
		isLoggingOut = true;
		try {
			await api.auth.logout(fetch);
		} catch (err) {
			console.error('Logout failed:', err);
			errorMessage = err instanceof Error ? err.message : 'Failed to disconnect Google account';
		} finally {
			isLoggingOut = false;
		}
	}
</script>

<div class="container mx-auto max-w-4xl p-4">
	<div class="mb-6">
		<h1 class="text-3xl font-bold tracking-tight">Settings</h1>
		<p class="text-muted-foreground">Manage your provider profile and business information</p>
	</div>

	{#if errorMessage}
		<ErrorAlert
			message={errorMessage}
			title="Error"
			showRetryButton={false}
			onRetry={() => {
				errorMessage = null;
			}}
		/>
	{:else if currentProvider}
		<ProviderForm
			provider={currentProvider}
			mode="edit"
			onSave={handleSave}
			onCancel={handleCancel}
		/>
	{:else}
		<Card.Root>
			<Card.Content class="flex flex-col items-center justify-center py-12">
				<AlertCircleIcon class="mb-4 h-12 w-12 text-muted-foreground" />
				<h3 class="mb-2 text-lg font-semibold">No Provider Selected</h3>
				<p class="mb-6 text-center text-sm text-muted-foreground">
					You need to create a provider profile to start creating invoices
				</p>
				<Button onclick={navigateToNewProvider}>Create Provider Profile</Button>
			</Card.Content>
		</Card.Root>
	{/if}

	<!-- Email Authentication Section (only show if OAuth2 configured) -->
	{#if $requiresOAuth}
		<Card.Root class="mt-8">
			<Card.Header>
				<Card.Title>Email Authentication</Card.Title>
				<Card.Description>
					Manage your Google account connection for sending emails
				</Card.Description>
			</Card.Header>
			<Card.Content>
				{#if $isAuthenticated}
					<div class="flex items-center justify-between">
						<div class="flex items-center gap-3">
							<div
								class="flex h-10 w-10 items-center justify-center rounded-full bg-green-100 dark:bg-green-900"
							>
								<CheckCircle2Icon class="h-5 w-5 text-green-600 dark:text-green-400" />
							</div>
							<div>
								<p class="text-sm font-medium">Connected</p>
								<p class="text-xs text-muted-foreground">{$currentUserEmail}</p>
							</div>
						</div>
						<Button variant="outline" onclick={handleLogout} disabled={isLoggingOut}>
							{isLoggingOut ? 'Disconnecting...' : 'Disconnect'}
						</Button>
					</div>
				{:else}
					<div class="flex items-center justify-between">
						<div class="flex items-center gap-3">
							<div
								class="flex h-10 w-10 items-center justify-center rounded-full bg-yellow-100 dark:bg-yellow-900"
							>
								<AlertCircleIcon class="h-5 w-5 text-yellow-600 dark:text-yellow-400" />
							</div>
							<div>
								<p class="text-sm font-medium">Not Connected</p>
								<p class="text-xs text-muted-foreground">
									Connect your Google account to send emails
								</p>
							</div>
						</div>
						<Button onclick={handleLogin} disabled={isLoggingIn}>
							{isLoggingIn ? 'Connecting...' : 'Sign in with Google'}
						</Button>
					</div>
				{/if}
			</Card.Content>
		</Card.Root>
	{/if}

	<!-- Version Info -->
	<div class="mt-8 text-center text-sm text-muted-foreground">
		Version: {version}
	</div>
</div>
