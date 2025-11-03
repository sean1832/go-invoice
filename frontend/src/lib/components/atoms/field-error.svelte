<!--
	Field Error Atom
	
	Simple component for displaying validation error messages below form fields.
	
	Props:
	- error: string | null | undefined - Error message to display
	- class?: string - Additional CSS classes
	
	Usage:
	<FieldError error={emailError} />
	{#if errors.length > 0}
		<FieldError error={errors[0]} />
	{/if}
-->
<script lang="ts">
	import { cn } from '@/utils';

	interface Props {
		error?: string | null;
		class?: string;
	}

	let { error, class: customClass = '' }: Props = $props();
</script>

{#if error}
	<p class={cn('text-sm text-destructive', customClass)}>
		{error}
	</p>
{/if}
