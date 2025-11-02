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
		<div class="flex flex-col sm:flex-row sm:justify-between sm:items-start mb-6 sm:mb-8 gap-6">
			<!-- Left: Invoice Info -->
			<div>
				<h1 class="text-3xl sm:text-4xl font-bold text-foreground mb-2">INVOICE</h1>
				<div class="space-y-1 text-sm sm:text-base">
					<p class="text-muted-foreground">
						Invoice Number: <span class="font-semibold text-foreground">{invoice.id}</span>
					</p>
					<p class="text-muted-foreground">
						Date: <span class="font-semibold text-foreground">{formatDateShort(invoice.date)}</span>
					</p>
					<p class="text-muted-foreground">
						Due Date: <span class="font-semibold text-foreground">{formatDateShort(invoice.due)}</span>
					</p>
				</div>
			</div>

			<!-- Right: Provider Info -->
			<div class="sm:text-right">
				<h2 class="text-xl sm:text-2xl font-bold text-foreground mb-2">{invoice.provider.name}</h2>
				<div class="space-y-0.5 text-sm sm:text-base text-muted-foreground">
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
			<h3 class="text-base sm:text-lg font-semibold text-foreground mb-3">Bill To:</h3>
			<div class="bg-muted p-3 sm:p-4 rounded-lg">
				<p class="font-semibold text-foreground">{invoice.client.name}</p>
				<div class="space-y-0.5 text-sm sm:text-base text-muted-foreground mt-1">
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
		<div class="hidden lg:block mb-8 overflow-hidden rounded-lg border-2 border-border">
			<table class="w-full">
				<thead>
					<tr class="bg-primary text-primary-foreground">
						<th class="py-4 px-6 text-left font-semibold w-36">Date</th>
						<th class="py-4 px-6 text-left font-semibold">Description</th>
						<th class="py-4 px-6 text-center font-semibold w-28">Qty</th>
						<th class="py-4 px-6 text-right font-semibold w-36">Unit Price</th>
						<th class="py-4 px-6 text-right font-semibold w-36">Amount</th>
					</tr>
				</thead>
				<tbody class="bg-card">
					{#each invoice.items as item}
						<tr class="border-b-2 border-border hover:bg-muted/50 transition-colors">
							<td class="py-4 px-6 text-left text-muted-foreground">{formatDateShort(item.date)}</td>
							<td class="py-4 px-6">
								<p class="font-semibold text-foreground">{item.description}</p>
								{#if item.descriptionDetail}
									<p class="text-muted-foreground text-sm mt-1">{item.descriptionDetail}</p>
								{/if}
							</td>
							<td class="py-4 px-6 text-center text-muted-foreground">{item.quantity}</td>
							<td class="py-4 px-6 text-right text-muted-foreground">{formatCurrency(item.unitPrice)}</td>
							<td class="py-4 px-6 text-right font-semibold text-foreground">
								{formatCurrency(item.totalPrice)}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		<!-- Invoice Items - Mobile Card View -->
		<div class="lg:hidden mb-6 sm:mb-8 space-y-4">
			{#each invoice.items as item}
				<div class="border-2 border-border rounded-lg p-4 hover:border-primary/50 transition-colors bg-card">
					<!-- Item Header -->
					<div class="flex justify-between items-start mb-2">
						<div class="flex-1">
							<p class="font-semibold text-foreground mb-1">{item.description}</p>
							{#if item.descriptionDetail}
								<p class="text-muted-foreground text-sm mb-2">{item.descriptionDetail}</p>
							{/if}
							<p class="text-xs text-muted-foreground">{formatDateShort(item.date)}</p>
						</div>
					</div>

					<Separator class="my-3" />

					<!-- Item Details Grid -->
					<div class="grid grid-cols-2 gap-3 text-sm">
						<div>
							<p class="text-muted-foreground text-xs mb-0.5">Quantity</p>
							<p class="font-semibold text-foreground">{item.quantity}</p>
						</div>
						<div>
							<p class="text-muted-foreground text-xs mb-0.5">Unit Price</p>
							<p class="font-semibold text-foreground">{formatCurrency(item.unitPrice)}</p>
						</div>
					</div>

					<!-- Item Total -->
					<div class="mt-3 pt-3 border-t border-border">
						<div class="flex justify-between items-center">
							<p class="text-muted-foreground font-medium">Amount</p>
							<p class="text-lg font-bold text-foreground">{formatCurrency(item.totalPrice)}</p>
						</div>
					</div>
				</div>
			{/each}
		</div>

		<!-- Totals Section -->
		<div class="flex justify-end mb-6 sm:mb-8">
			<div class="w-full sm:w-72">
				<div class="flex justify-between py-2 border-b border-border text-sm sm:text-base">
					<span class="text-muted-foreground">Subtotal:</span>
					<span class="font-semibold text-foreground">{formatCurrency(invoice.pricing.subtotal)}</span>
				</div>
				<div class="flex justify-between py-2 border-b border-border text-sm sm:text-base">
					<span class="text-muted-foreground">Tax ({invoice.pricing.taxRate}%):</span>
					<span class="font-semibold text-foreground">{formatCurrency(invoice.pricing.tax)}</span>
				</div>
				<div
					class="flex justify-between py-3 bg-primary text-primary-foreground px-4 rounded-lg mt-2 text-base sm:text-lg"
				>
					<span class="font-bold">Total:</span>
					<span class="font-bold">{formatCurrency(invoice.pricing.total)}</span>
				</div>
			</div>
		</div>

		<!-- Payment Information -->
		<div class="border-t border-border pt-4 sm:pt-6">
			<h3 class="text-base sm:text-lg font-semibold text-foreground mb-3">Payment Information</h3>
			<div class="grid grid-cols-1 sm:grid-cols-2 gap-3 text-sm">
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
		<div class="mt-6 sm:mt-8 border-t border-border pt-4 sm:pt-6">
			<p class="text-xs sm:text-sm text-muted-foreground">
				Payment is due within 30 days. Please include the invoice number with your payment. Thank
				you for your business!
			</p>
		</div>
	</Card.Content>
</Card.Root>
