<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '@/components/ui/input';
	import { Textarea } from '@/components/ui/textarea';
	import AttachmentIcon from '@lucide/svelte/icons/paperclip';
	import Label from '@/components/ui/label/label.svelte';
	import Button from '@/components/ui/button/button.svelte';
	import Spinner from '@/components/atoms/spinner.svelte';
	import type { EmailConfig, EmailContent } from '@/types/invoice';
	interface Props {
		children: () => any;
		onSendEmail?: (data: EmailConfig) => Promise<void> | void;
		templateData: EmailConfig;
		isSending: boolean;
	}
	let { children, onSendEmail: onSubmit, templateData, isSending }: Props = $props();

	// dialog open state
	let open = $state(false);
	let wasSending = $state(false);

	// derive local state from templateData - convert array to comma-separated string
	let emailTo = $state(templateData.to.join(', '));
	let emailSubject = $state(templateData.subject);
	let emailBody = $state(templateData.body);

	async function handleSubmit() {
		if (onSubmit) {
			// parse emailTo with comma separated values and trim spaces
			const recipients = emailTo
				.split(',')
				.map((email) => email.trim())
				.filter((email) => email.length > 0);

			await onSubmit({
				to: recipients,
				subject: emailSubject,
				body: emailBody
			});
		}
	}

	// Track sending state and close dialog after successful submission
	$effect(() => {
		if (wasSending && !isSending) {
			// Email was sending, now it's done - close the dialog
			open = false;
		}
		wasSending = isSending;
	});

	// Reset form when dialog closes
	$effect(() => {
		if (!open) {
			emailTo = templateData.to.join(', ');
			emailSubject = templateData.subject;
			emailBody = templateData.body;
		}
	});
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger>
		{@render children()}
	</Dialog.Trigger>

	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Send Email</Dialog.Title>
		</Dialog.Header>
		<div class="flex flex-col gap-2">
			<div class="relative">
				<Label for="email_to">To</Label>
				<Input
					id="email_to"
					placeholder="example@email.com"
					bind:value={emailTo}
					class="mt-2"
					disabled={isSending}
				/>
			</div>
			<div class="relative mt-8">
				<Label for="email_subject">Subject</Label>
				<Input id="email_subject" class="mt-2" bind:value={emailSubject} disabled={isSending} />
			</div>
			<div class="relative">
				<Label for="email_body">Body</Label>
				<Textarea id="email_body" class="mt-2" bind:value={emailBody} disabled={isSending} />
			</div>
			<div>
				<Button variant="ghost" disabled={isSending}>
					<AttachmentIcon class="mt-2 mr-1 inline h-4 w-4" />
				</Button>
			</div>
		</div>
		<Dialog.Footer>
			<div class="flex gap-2">
				<Dialog.Close>
					<Button variant="outline" disabled={isSending}>Cancel</Button>
				</Dialog.Close>
				<Button type="submit" onclick={handleSubmit} disabled={isSending}>
					{#if isSending}
						<Spinner class="mr-2 h-4 w-4" size={16} />
					{/if}
					{isSending ? 'Sending...' : 'Send'}
				</Button>
			</div>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
