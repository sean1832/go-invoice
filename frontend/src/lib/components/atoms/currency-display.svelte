<!--
	Currency Display Atom
	
	Simple, reusable component for displaying formatted currency values.
	
	Props:
	- amount: number - The amount to display
	- currency?: string - Currency code (default: AUD)
	- class?: string - Additional CSS classes
	
	Usage:
	<CurrencyDisplay amount={123.45} />
	<CurrencyDisplay amount={total} class="text-lg font-bold" />
-->
<script lang="ts">
	import { formatCurrency } from '@/utils/formatters';
	import { cn } from '@/utils';

	interface Props {
		amount: number;
		currency?: string;
		class?: string;
	}

	let { amount, currency = 'AUD', class: customClass = '' }: Props = $props();
</script>

<span class={cn('tabular-nums', customClass)}>
	{formatCurrency(amount, currency)}
</span>
