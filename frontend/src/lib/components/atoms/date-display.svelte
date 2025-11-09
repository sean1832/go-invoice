<!--
	Date Display Atom
	
	Simple, reusable component for displaying formatted dates.
	
	Props:
	- date: string - ISO date string (YYYY-MM-DD)
	- format?: 'long' | 'short' | 'numeric' - Display format
	- class?: string - Additional CSS classes
	
	Usage:
	<DateDisplay date="2025-11-03" />
	<DateDisplay date={invoice.date} format="short" />
-->
<script lang="ts">
	import { formatDate, formatDateShort, formatDateNumeric } from '@/helpers';
	import { cn } from '@/utils';

	interface Props {
		date: string;
		format?: 'long' | 'short' | 'numeric';
		class?: string;
	}

	let { date, format = 'long', class: customClass = '' }: Props = $props();

	const formattedDate = $derived(() => {
		switch (format) {
			case 'short':
				return formatDateShort(date);
			case 'numeric':
				return formatDateNumeric(date);
			default:
				return formatDate(date);
		}
	});
</script>

<span class={cn(customClass)}>
	{formattedDate()}
</span>
