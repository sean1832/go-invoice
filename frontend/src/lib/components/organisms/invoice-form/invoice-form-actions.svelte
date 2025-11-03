<!--
	Invoice Form Actions Organism
	
	Action buttons for invoice forms (Cancel, Save Draft, Send).
	
	Props:
	- onCancel: () => void - Callback for cancel button
	- onSaveDraft: () => void - Callback for save draft button
	- onSend: () => void - Callback for send button
	- isSaving?: boolean - Loading state for buttons
	- class?: string - Additional CSS classes
	
	Usage:
	<InvoiceFormActions
		onCancel={handleCancel}
		onSaveDraft={handleSaveDraft}
		onSend={handleSend}
	/>
-->
<script lang="ts">
	import { cn } from '@/utils';
	import Button from '@/components/ui/button/button.svelte';
	import SaveIcon from '@lucide/svelte/icons/save';
	import SendIcon from '@lucide/svelte/icons/send';

	interface Props {
		onCancel: () => void;
		onSaveDraft: () => void;
		onSend: () => void;
		isSaving?: boolean;
		class?: string;
	}

	let {
		onCancel,
		onSaveDraft,
		onSend,
		isSaving = false,
		class: customClass = ''
	}: Props = $props();
</script>

<div class={cn('flex justify-between', customClass)}>
	<Button variant="outline" onclick={onCancel} disabled={isSaving}>Cancel</Button>

	<div class="flex gap-2">
		<Button variant="outline" onclick={onSaveDraft} disabled={isSaving}>
			<SaveIcon class="mr-2 h-4 w-4" />
			{isSaving ? 'Saving...' : 'Save as Draft'}
		</Button>
		<Button variant="default" onclick={onSend} disabled={isSaving}>
			<SendIcon class="mr-2 h-4 w-4" />
			{isSaving ? 'Sending...' : 'Save & Send'}
		</Button>
	</div>
</div>
