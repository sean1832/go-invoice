<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import Button from '@/components/ui/button/button.svelte';

	interface Props {
		/** Controls the open/closed state of the dialog */
		open?: boolean;
		/** Dialog title */
		title?: string;
		/** Dialog description/message */
		description?: string;
		/** Text for the confirm button */
		confirmText?: string;
		/** Text for the cancel button */
		cancelText?: string;
		/** Variant for the confirm button */
		confirmVariant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link';
		/** Callback when confirm is clicked */
		onConfirm: () => void;
		/** Callback when cancel is clicked (optional) */
		onCancel?: () => void;
	}

	let {
		open = $bindable(false),
		title = 'Confirm Action',
		description = 'Are you sure you want to proceed?',
		confirmText = 'Confirm',
		cancelText = 'Cancel',
		confirmVariant = 'destructive',
		onConfirm,
		onCancel
	}: Props = $props();

	function handleConfirm() {
		onConfirm();
		open = false;
	}

	function handleCancel() {
		if (onCancel) {
			onCancel();
		}
		open = false;
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>{title}</Dialog.Title>
			<Dialog.Description>
				{@html description}
			</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button variant="outline" onclick={handleCancel}>{cancelText}</Button>
			<Button variant={confirmVariant} onclick={handleConfirm}>{confirmText}</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
