<script lang="ts">
	import * as Alert from '$lib/components/ui/alert/index.js';
	import Button from '@/components/ui/button/button.svelte';
	import { cn } from '@/utils';
	import AlertCircleIcon from '@lucide/svelte/icons/alert-circle';
	import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';

	interface Props {
		/**
		 * The error message to display
		 */
		message: string | null | undefined;

		type?: 'error' | 'warning' | 'info' | undefined;

		/**
		 * Alert title
		 */
		title?: string;
		/**
		 * Whether to show a back button
		 */
		showBackButton?: boolean;
		/**
		 * Custom back button text (default: "Back")
		 */
		backButtonText?: string;
		/**
		 * Custom back button handler. If not provided, uses window.history.back()
		 */
		onBack?: () => void;
		/**
		 * Whether to show a retry button
		 */
		showRetryButton?: boolean;
		/**
		 * Retry button text (default: "Retry")
		 */
		retryButtonText?: string;
		/**
		 * Retry button handler
		 */
		onRetry?: () => void;
		/**
		 * Additional CSS classes for the container
		 */
		class?: string;
	}

	let {
		message,
		type = 'error',
		title: titleProp,
		showBackButton = false,
		backButtonText = 'Back',
		onBack,
		showRetryButton = false,
		retryButtonText = 'Retry',
		onRetry,
		class: className = ''
	}: Props = $props();

	const title = $derived(
		titleProp ?? (type === 'error' ? 'Error' : type === 'warning' ? 'Warning' : 'Info')
	);

	function handleBack() {
		if (onBack) {
			onBack();
		} else {
			window.history.back();
		}
	}
</script>

{#if message}
	<div class={cn('mb-6', className)}>
		{#if showBackButton}
			<Button variant="ghost" size="sm" onclick={handleBack} class="mb-4">
				<ArrowLeftIcon class="mr-2 h-4 w-4" />
				{backButtonText}
			</Button>
		{/if}
		<Alert.Root
			variant={type === 'error' ? 'destructive' : type === 'warning' ? 'warn' : 'default'}
		>
			<AlertCircleIcon class="h-4 w-4" />
			<Alert.Title>{title}</Alert.Title>
			<Alert.Description>
				{message}
			</Alert.Description>
		</Alert.Root>
		{#if showRetryButton && onRetry}
			<div class="mt-4 flex justify-center">
				<Button onclick={onRetry} class="w-full cursor-pointer">
					{retryButtonText}
				</Button>
			</div>
		{/if}
	</div>
{/if}
