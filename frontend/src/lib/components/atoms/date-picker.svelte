<script lang="ts">
	import CalendarIcon from '@lucide/svelte/icons/calendar';
	import { DateFormatter, type DateValue, getLocalTimeZone } from '@internationalized/date';
	import { cn } from '$lib/utils.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { Calendar } from '$lib/components/ui/calendar/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';

	// Define the props for reusability
	type Props = {
		value: DateValue | undefined; // Use `bind:value` in the parent
		placeholder?: string;
		locale?: string;
		formatOptions?: Intl.DateTimeFormatOptions;
		class?: string; // To pass in custom classes (e.g., width)
	};

	// Use $props() to get props, providing defaults
	let {
		value = $bindable(),
		placeholder = 'Pick a date',
		locale = 'en-US',
		formatOptions = { dateStyle: 'long' },
		class: customClass = '' // Rename `class` to avoid conflict
	}: Props = $props();

	// Internal state for the popover content reference
	let contentRef = $state<HTMLElement | null>(null);
	let open = $state(false);

	// Use $derived for computed values that react to prop changes
	const df = $derived(new DateFormatter(locale, formatOptions));

	const displayText = $derived(value ? df.format(value.toDate(getLocalTimeZone())) : placeholder);

	// Close popover when a date is selected
	$effect(() => {
		if (value) {
			open = false;
		}
	});
</script>

<Popover.Root bind:open>
	<Popover.Trigger
		class={cn(
			buttonVariants({
				variant: 'outline',
				class: 'justify-start text-left font-normal' // Removed fixed width
			}),
			!value && 'text-muted-foreground',
			customClass // Apply any additional classes from the parent
		)}
	>
		<CalendarIcon class="mr-2 h-4 w-4" />
		{displayText}
	</Popover.Trigger>
	<Popover.Content bind:ref={contentRef} class="w-auto p-0">
		<Calendar type="single" bind:value />
	</Popover.Content>
</Popover.Root>
