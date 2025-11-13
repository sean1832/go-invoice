<script lang="ts">
	import Button from '@/components/ui/button/button.svelte';
	import * as Card from '@/components/ui/card';
	import { isAuthenticated, currentUserEmail } from '@/stores';
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import CheckCircle2Icon from '@lucide/svelte/icons/check-circle-2';
	import { api } from '@/services';

	interface Props {
		/**
		 * Title of the card
		 * @default "Email Authentication"
		 */
		title?: string;
		/**
		 * Description text below the title
		 * @default "Manage your Google account connection for sending emails"
		 */
		description?: string;
		/**
		 * Text shown when connected
		 * @default "Connected"
		 */
		connectedText?: string;
		/**
		 * Text shown when not connected
		 * @default "Not Connected"
		 */
		notConnectedText?: string;
		/**
		 * Helper text shown when not connected
		 * @default "Connect your Google account to send emails"
		 */
		notConnectedHelper?: string;
		/**
		 * Text for the login button
		 * @default "Sign in with Google"
		 */
		loginButtonText?: string;
		/**
		 * Text for the logout button
		 * @default "Disconnect"
		 */
		logoutButtonText?: string;
		/**
		 * Whether to show the card wrapper
		 * @default true
		 */
		showCard?: boolean;
		/**
		 * Custom class for the card/container
		 */
		class?: string;
		/**
		 * Callback when login succeeds
		 */
		onLoginSuccess?: () => void;
		/**
		 * Callback when login fails
		 */
		onLoginError?: (error: Error) => void;
		/**
		 * Callback when logout succeeds
		 */
		onLogoutSuccess?: () => void;
		/**
		 * Callback when logout fails
		 */
		onLogoutError?: (error: Error) => void;
	}

	let {
		title = 'Email Authentication',
		description = 'Manage your Google account connection for sending emails',
		connectedText = 'Connected',
		notConnectedText = 'Not Connected',
		notConnectedHelper = 'Connect your Google account to send emails',
		loginButtonText = 'Sign in with Google',
		logoutButtonText = 'Disconnect',
		showCard = true,
		class: className = '',
		onLoginSuccess,
		onLoginError,
		onLogoutSuccess,
		onLogoutError
	}: Props = $props();

	let isLoggingIn = $state(false);
	let isLoggingOut = $state(false);

	async function handleLogin() {
		isLoggingIn = true;
		try {
			await api.auth.loginWithGoogle();
			onLoginSuccess?.();
		} catch (err) {
			console.error('Login failed:', err);
			const error = err instanceof Error ? err : new Error('Failed to connect Google account');
			onLoginError?.(error);
		} finally {
			isLoggingIn = false;
		}
	}

	async function handleLogout() {
		isLoggingOut = true;
		try {
			await api.auth.logout(fetch);
			onLogoutSuccess?.();
		} catch (err) {
			console.error('Logout failed:', err);
			const error = err instanceof Error ? err : new Error('Failed to disconnect Google account');
			onLogoutError?.(error);
		} finally {
			isLoggingOut = false;
		}
	}
</script>

{#if showCard}
	<Card.Root class={className}>
		<Card.Header>
			<Card.Title>{title}</Card.Title>
			<Card.Description>{description}</Card.Description>
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
							<p class="text-sm font-medium">{connectedText}</p>
							<p class="text-xs text-muted-foreground">{$currentUserEmail}</p>
						</div>
					</div>
					<Button variant="outline" onclick={handleLogout} disabled={isLoggingOut}>
						{isLoggingOut ? 'Disconnecting...' : logoutButtonText}
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
							<p class="text-sm font-medium">{notConnectedText}</p>
							<p class="text-xs text-muted-foreground">{notConnectedHelper}</p>
						</div>
					</div>
					<Button onclick={handleLogin} disabled={isLoggingIn}>
						{isLoggingIn ? 'Connecting...' : loginButtonText}
					</Button>
				</div>
			{/if}
		</Card.Content>
	</Card.Root>
{:else}
	<div class={className}>
		{#if $isAuthenticated}
			<div class="flex items-center justify-between">
				<div class="flex items-center gap-3">
					<div
						class="flex h-10 w-10 items-center justify-center rounded-full bg-green-100 dark:bg-green-900"
					>
						<CheckCircle2Icon class="h-5 w-5 text-green-600 dark:text-green-400" />
					</div>
					<div>
						<p class="text-sm font-medium">{connectedText}</p>
						<p class="text-xs text-muted-foreground">{$currentUserEmail}</p>
					</div>
				</div>
				<Button variant="outline" onclick={handleLogout} disabled={isLoggingOut}>
					{isLoggingOut ? 'Disconnecting...' : logoutButtonText}
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
						<p class="text-sm font-medium">{notConnectedText}</p>
						<p class="text-xs text-muted-foreground">{notConnectedHelper}</p>
					</div>
				</div>
				<Button onclick={handleLogin} disabled={isLoggingIn}>
					{isLoggingIn ? 'Connecting...' : loginButtonText}
				</Button>
			</div>
		{/if}
	</div>
{/if}
