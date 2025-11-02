<!--
	InvoiceCard Component
	
	A responsive, modular invoice display component.
	
	Props:
	- invoice: Invoice object containing all invoice data
	
	Responsive Design:
	- Desktop (lg+): Shows items in a table format with columns for Date, Description, Qty, Unit Price, Amount
	- Mobile/Tablet (<lg): Shows items as individual cards with vertical layout
	- All text sizes and spacing adapt to screen size using Tailwind breakpoints (sm:, lg:)
	- Action buttons show icon-only on mobile, icon+text on desktop
	
	Usage:
	<InvoiceCard {invoice} />
-->
<script lang="ts">
	import * as Card from '@/components/ui/card';
	import Separator from '@/components/ui/separator/separator.svelte';
	import type { Invoice } from '@/types/invoice';

	interface Props {
		invoice: Invoice;
	}

	let { invoice }: Props = $props();

	// Helper functions
	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-AU', {
			style: 'currency',
			currency: 'AUD'
		}).format(amount);
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}

	function formatDateShort(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	}
</script>

<Card.Root>
	<Card.Content class="p-4 sm:p-8">
		<!-- Header Section -->
		<div class="mb-6 flex flex-col gap-6 sm:mb-8 sm:flex-row sm:items-start sm:justify-between">
			<!-- Left: Invoice Info -->
			<div>
				<h1 class="mb-2 text-3xl font-bold text-foreground sm:text-4xl">INVOICE</h1>
				<div class="space-y-1 text-sm sm:text-base">
					<p class="text-muted-foreground">
						Invoice Number: <span class="font-semibold text-foreground">{invoice.id}</span>
					</p>
					<p class="text-muted-foreground">
						Date: <span class="font-semibold text-foreground">{formatDateShort(invoice.date)}</span>
					</p>
					<p class="text-muted-foreground">
						Due Date: <span class="font-semibold text-foreground"
							>{formatDateShort(invoice.due)}</span
						>
					</p>
				</div>
			</div>

			<!-- Right: Provider Info -->
			<div class="sm:text-right">
				<h2 class="mb-2 text-xl font-bold text-foreground sm:text-2xl">{invoice.provider.name}</h2>
				<div class="space-y-0.5 text-sm text-muted-foreground sm:text-base">
					{#if invoice.provider.address}
						<p class="whitespace-pre-line">{invoice.provider.address}</p>
					{/if}
					{#if invoice.provider.abn}
						<p>ABN: {invoice.provider.abn}</p>
					{/if}
					{#if invoice.provider.phone}
						<p>Phone: {invoice.provider.phone}</p>
					{/if}
					{#if invoice.provider.email}
						<p>Email: {invoice.provider.email}</p>
					{/if}
				</div>
			</div>
		</div>

		<!-- Bill To Section -->
		<div class="mb-6 sm:mb-8">
			<h3 class="mb-3 text-base font-semibold text-foreground sm:text-lg">Bill To:</h3>
			<div class="rounded-lg bg-muted p-3 sm:p-4">
				<p class="font-semibold text-foreground">{invoice.client.name}</p>
				<div class="mt-1 space-y-0.5 text-sm text-muted-foreground sm:text-base">
					{#if invoice.client.address}
						<p class="whitespace-pre-line">{invoice.client.address}</p>
					{/if}
					{#if invoice.client.abn}
						<p>ABN: {invoice.client.abn}</p>
					{/if}
					{#if invoice.client.phone}
						<p>Phone: {invoice.client.phone}</p>
					{/if}
					{#if invoice.client.email}
						<p>Email: {invoice.client.email}</p>
					{/if}
				</div>
			</div>
		</div>

		<!-- Invoice Items - Desktop Table View -->
		<div class="mb-8 hidden overflow-hidden rounded-lg border-2 border-border lg:block">
			<table class="w-full">
				<thead>
					<tr class="bg-primary text-primary-foreground">
						<th class="w-36 px-6 py-4 text-left font-semibold">Date</th>
						<th class="px-6 py-4 text-left font-semibold">Description</th>
						<th class="w-28 px-6 py-4 text-center font-semibold">Qty</th>
						<th class="w-36 px-6 py-4 text-right font-semibold">Unit Price</th>
						<th class="w-36 px-6 py-4 text-right font-semibold">Amount</th>
					</tr>
				</thead>
				<tbody class="bg-card">
					{#each invoice.items as item}
						<tr class="border-b-2 border-border transition-colors hover:bg-muted/50">
							<td class="px-6 py-4 text-left text-muted-foreground">{formatDateShort(item.date)}</td
							>
							<td class="px-6 py-4">
								<p class="font-semibold text-foreground">{item.description}</p>
								{#if item.descriptionDetail}
									<p class="mt-1 text-sm text-muted-foreground">{item.descriptionDetail}</p>
								{/if}
							</td>
							<td class="px-6 py-4 text-center text-muted-foreground">{item.quantity}</td>
							<td class="px-6 py-4 text-right text-muted-foreground"
								>{formatCurrency(item.unitPrice)}</td
							>
							<td class="px-6 py-4 text-right font-semibold text-foreground">
								{formatCurrency(item.totalPrice)}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<!-- Invoice Items - Mobile Card View -->
		<div class="mb-6 space-y-4 sm:mb-8 lg:hidden">
			{#each invoice.items as item}
				<div
					class="rounded-lg border-2 border-border bg-card p-4 transition-colors hover:border-primary/50"
				>
					<!-- Item Header -->
					<div class="mb-2 flex items-start justify-between">
						<div class="flex-1">
							<p class="mb-1 font-semibold text-foreground">{item.description}</p>
							{#if item.descriptionDetail}
								<p class="mb-2 text-sm text-muted-foreground">{item.descriptionDetail}</p>
							{/if}
							<p class="text-xs text-muted-foreground">{formatDateShort(item.date)}</p>
						</div>
					</div>

					<Separator class="my-3" />

					<!-- Item Details Grid -->
					<div class="grid grid-cols-2 gap-3 text-sm">
						<div>
							<p class="mb-0.5 text-xs text-muted-foreground">Quantity</p>
							<p class="font-semibold text-foreground">{item.quantity}</p>
						</div>
						<div>
							<p class="mb-0.5 text-xs text-muted-foreground">Unit Price</p>
							<p class="font-semibold text-foreground">{formatCurrency(item.unitPrice)}</p>
						</div>
					</div>

					<!-- Item Total -->
					<div class="mt-3 border-t border-border pt-3">
						<div class="flex items-center justify-between">
							<p class="font-medium text-muted-foreground">Amount</p>
							<p class="text-lg font-bold text-foreground">{formatCurrency(item.totalPrice)}</p>
						</div>
					</div>
				</div>
			{/each}
		</div>

		<!-- Totals Section -->
		<div class="mb-6 flex justify-end sm:mb-8">
			<div class="w-full sm:w-72">
				<div class="flex justify-between border-b border-border py-2 text-sm sm:text-base">
					<span class="text-muted-foreground">Subtotal:</span>
					<span class="font-semibold text-foreground"
						>{formatCurrency(invoice.pricing.subtotal)}</span
					>
				</div>
				<div class="flex justify-between border-b border-border py-2 text-sm sm:text-base">
					<span class="text-muted-foreground">Tax ({invoice.pricing.taxRate}%):</span>
					<span class="font-semibold text-foreground">{formatCurrency(invoice.pricing.tax)}</span>
				</div>
				<div
					class="mt-2 flex justify-between rounded-lg bg-primary px-4 py-3 text-base text-primary-foreground sm:text-lg"
				>
					<span class="font-bold">Total:</span>
					<span class="font-bold">{formatCurrency(invoice.pricing.total)}</span>
				</div>
			</div>
		</div>

		<!-- Payment Information -->
		<div class="border-t border-border pt-4 sm:pt-6">
			<h3 class="mb-3 text-base font-semibold text-foreground sm:text-lg">Payment Information</h3>
			<div class="grid grid-cols-1 gap-3 text-sm sm:grid-cols-2">
				<div class="space-y-2">
					<p class="text-muted-foreground">
						<span class="font-semibold text-foreground">Payment Method:</span>
						{invoice.payment.method}
					</p>
					<p class="text-muted-foreground">
						<span class="font-semibold text-foreground">Account Name:</span>
						{invoice.payment.accountName}
					</p>
				</div>
				<div class="space-y-2">
					<p class="text-muted-foreground">
						<span class="font-semibold text-foreground">BSB:</span>
						{invoice.payment.bsb}
					</p>
					<p class="text-muted-foreground">
						<span class="font-semibold text-foreground">Account Number:</span>
						{invoice.payment.accountNumber}
					</p>
				</div>
			</div>
		</div>

		<!-- Notes/Terms -->
		<div class="mt-6 border-t border-border pt-4 sm:mt-8 sm:pt-6">
			<p class="text-xs text-muted-foreground sm:text-sm">
				Payment is due within 30 days. Please include the invoice number with your payment. Thank
				you for your business!
			</p>
		</div>
	</Card.Content>
</Card.Root>
