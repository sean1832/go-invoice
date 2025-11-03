<!--
	Status Badge Atom
	
	Simple badge component for displaying invoice status with appropriate styling.
	
	Props:
	- status: InvoiceStatus - The status to display ('draft' | 'send')
	- class?: string - Additional CSS classes
	
	Usage:
	<StatusBadge status="draft" />
	<StatusBadge status={invoice.status} />
-->
<script lang="ts">
	import type { InvoiceStatus } from '@/types/invoice';
	import { cn } from '@/utils';
	import Badge from '@/components/ui/badge/badge.svelte';

	interface Props {
		status: InvoiceStatus;
		class?: string;
	}

	let { status, class: customClass = '' }: Props = $props();

	const statusConfig = $derived(() => {
		switch (status) {
			case 'draft':
				return {
					label: 'Draft',
					variant: 'secondary' as const,
					class: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-100'
				};
			case 'send':
				return {
					label: 'Sent',
					variant: 'default' as const,
					class: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100'
				};
			default:
				return {
					label: String(status),
					variant: 'outline' as const,
					class: ''
				};
		}
	});
</script>

<Badge variant={statusConfig().variant} class={cn(statusConfig().class, customClass)}>
	{statusConfig().label}
</Badge>
