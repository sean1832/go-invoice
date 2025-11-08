<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '@/components/ui/input';
	import { Textarea } from '@/components/ui/textarea';
	import AttachmentIcon from '@lucide/svelte/icons/paperclip';
	import Label from '@/components/ui/label/label.svelte';
	import Button from '@/components/ui/button/button.svelte';
	interface Props {
		children: () => any;
		onSubmit?: (data: any) => void;
	}
	let { children, onSubmit }: Props = $props();

	function handleSubmit() {}
</script>

<Dialog.Root>
	<Dialog.Trigger>
		{@render children()}
	</Dialog.Trigger>

	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Send Email</Dialog.Title>
		</Dialog.Header>
		<div class="flex flex-col gap-2">
			<div class="relative">
				<Label for="email_from">From</Label>
				<Input id="email_from" value="user@example.com" disabled class="mt-2" />
			</div>
			<div class="relative">
				<Label for="email_to">To</Label>
				<Input id="email_to" placeholder="example@email.com" class="mt-2" />
			</div>
			<div class="relative mt-8">
				<Label for="email_subject">Subject</Label>
				<Input id="email_subject" class="mt-2" />
			</div>
			<div class="relative">
				<Label for="email_body">Body</Label>
				<Textarea id="email_body" class="mt-2" />
			</div>
			<div>
				<Button variant="ghost">
					<AttachmentIcon class="mt-2 mr-1 inline h-4 w-4" />
				</Button>
			</div>
		</div>
		<Dialog.Footer>
			<div class="flex gap-2">
				<Dialog.Close>
					<Button variant="outline">Cancel</Button>
				</Dialog.Close>
				<Button type="submit" onclick={handleSubmit}>Send</Button>
			</div>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
