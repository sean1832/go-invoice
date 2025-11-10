<!--
	Payment Info Card Molecule
	
	Displays payment information (bank details) for invoices.
	
	Props:
	- payment: PaymentInfo - Payment information to display
	- class?: string - Additional CSS classes
	
	Usage:
	<PaymentInfoCard payment={invoice.payment} />
-->
<script lang="ts">
	import { formatBSB } from '@/helpers';
	import type { PaymentInfo } from '@/types/invoice';
	import { cn } from '@/utils';

	interface Props {
		payment: PaymentInfo;
		class?: string;
	}

	let { payment, class: customClass = '' }: Props = $props();
</script>

<div class={cn('space-y-3', customClass)}>
	<h3 class="text-base font-semibold text-foreground sm:text-lg">Payment Information</h3>

	<div class="grid grid-cols-1 gap-3 text-sm sm:grid-cols-2">
		<div class="space-y-2">
			<p class="text-muted-foreground">
				<span class="font-semibold text-foreground">Payment Method:</span>
				{payment.method}
			</p>
			<p class="text-muted-foreground">
				<span class="font-semibold text-foreground">Account Name:</span>
				{payment.account_name}
			</p>
		</div>
		<div class="space-y-2">
			<p class="text-muted-foreground">
				<span class="font-semibold text-foreground">BSB:</span>
				{formatBSB(payment.bsb)}
			</p>
			<p class="text-muted-foreground">
				<span class="font-semibold text-foreground">Account Number:</span>
				{payment.account_number}
			</p>
		</div>
	</div>
</div>
