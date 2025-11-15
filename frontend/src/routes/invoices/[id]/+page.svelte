<script lang="ts">
	import { Badge } from '@/components/ui/badge';
	import Button from '@/components/ui/button/button.svelte';
	import { InvoiceDisplay } from '@/components/organisms/invoice-display';
	import type { EmailConfig, EmailTemplate, Invoice } from '@/types/invoice';
	import EditIcon from '@lucide/svelte/icons/pencil';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import SendIcon from '@lucide/svelte/icons/send';
	import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';
	import { api } from '@/services';
	import EmailDialog from '@/components/molecules/email-dialog.svelte';
	import { formatEmailTemplate, validateEmailConfig } from '$lib/helpers';
	import { toast } from 'svelte-sonner';
	import Spinner from '@/components/atoms/spinner.svelte';
	import { authMethod, invoices } from '@/stores';

	interface Props {
		data: {
			invoice: Invoice | null;
			emailConfig: EmailTemplate | null;
			error?: string;
		};
	}

	let { data }: Props = $props();
	let invoice = $state(data.invoice as Invoice);
	let emailData = $derived(data.emailConfig as EmailTemplate);
	let error = $derived(data.error);

	// Sync invoice state with data changes
	$effect(() => {
		if (data.invoice) {
			invoice = data.invoice as Invoice;
		}
	});

	// Format email template with invoice data
	let formattedEmail = $derived.by(() => {
		if (!emailData || !invoice) return { to: [], subject: '', body: '' } as EmailConfig;

		const pattern = {
			INVOICE_ID: invoice.id,
			CLIENT_NAME: invoice.client.name,
			PROVIDER_NAME: invoice.provider.name,
			PROVIDER_EMAIL: invoice.provider.email || '',
			SERVICE_TYPE: invoice.items[0].description || ''
		};

		// Parse email_target as comma-separated list if it exists
		const recipients = invoice.email_target
			? invoice.email_target.split(',').map((email) => email.trim())
			: [];

		return {
			to: recipients,
			subject: formatEmailTemplate(emailData.subject, pattern),
			body: formatEmailTemplate(emailData.body, pattern)
		} as EmailConfig;
	});

	// Helper functions
	function getStatusVariant(status: Invoice['status']): 'default' | 'secondary' {
		return status === 'send' ? 'default' : 'secondary';
	}

	function getStatusLabel(status: Invoice['status']): string {
		return status.charAt(0).toUpperCase() + status.slice(1);
	}

	// Action handlers
	function goBack() {
		// window.history.back();
		window.location.href = '/';
	}

	function editInvoice() {
		window.location.href = `/invoices/${invoice.id}/edit`;
	}

	let downloadError = $state<string | null>(null);
	let isDownloading = $state<boolean>(false);

	async function downloadInvoice() {
		try {
			isDownloading = true;
			const blob = await api.invoices.downloadPdf(fetch, invoice.id);

			// Create a temporary link to trigger the download
			const link = document.createElement('a');
			const url = URL.createObjectURL(blob);

			link.href = url;
			link.download = `${invoice.id}.pdf`;

			document.body.appendChild(link);
			link.click();

			// Clean up
			document.body.removeChild(link);
			URL.revokeObjectURL(url);
		} catch (error) {
			console.error('Error downloading invoice PDF:', error);
			downloadError =
				error instanceof Error ? error.message : 'Failed to download PDF. Please try again.';
		} finally {
			downloadError = null;
			isDownloading = false;
		}
	}

	let isSending = $state(false);

	async function onSendEmail(emailConfig: EmailConfig) {
		// validate first
		const validation = validateEmailConfig(emailConfig);
		if (!validation.isValid) {
			const errorMessage =
				validation.errors.length === 1
					? validation.errors[0]
					: `${validation.errors.length} errors: ${validation.errors.join(' • ')}`;
			toast.error('Email validation failed', {
				description: errorMessage
			});
			return;
		}

		try {
			isSending = true;
			await api.invoices.sendInvoiceEmail(fetch, invoice.id, emailConfig);

			// Refetch the invoice to get updated status from backend
			const updatedInvoice = await api.invoices.getInvoice(fetch, invoice.id);
			invoice = updatedInvoice;

			// Show success message
			toast.success('Invoice email sent successfully');
		} catch (error) {
			console.error('Error sending invoice email:', error);
			toast.error('Failed to send invoice email', {
				description: error instanceof Error ? error.message : 'Please try again.'
			});
		} finally {
			isSending = false;
		}
	}

	let isStatusUpdating = $state(false);
	async function onStatusBadgeClick() {
		// send back invoice with toggled status
		const currentInvoice = { ...invoice }; // shallow copy
		if (currentInvoice.status === 'draft') {
			currentInvoice.status = 'send';
		} else {
			currentInvoice.status = 'draft';
		}

		try {
			isStatusUpdating = true;

			const updatedInvoice = await api.invoices.updateInvoice(fetch, invoice.id, currentInvoice);
			invoice = updatedInvoice;
			toast.success('Invoice status updated.');
		} catch (error) {
			console.error('Error updating invoice status: ', error);
			toast.error('Failed to update invoice status', {
				description: error instanceof Error ? error.message : 'Please try again.'
			});
		} finally {
			isStatusUpdating = false;
		}
	}
</script>

<div class="container mx-auto max-w-5xl p-4">
	{#if error || !invoice}
		<ErrorAlert message={error || 'Invoice not found'} showBackButton={true} />
	{:else if downloadError}
		<ErrorAlert message={downloadError} showBackButton={true} />
	{:else}
		<!-- Invoice Display Component -->
		<InvoiceDisplay {invoice} />

		<!-- Action Buttons Bar -->
		<div class="mt-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
			<Button
				id="back"
				variant="ghost"
				size="sm"
				onclick={goBack}
				class="self-start"
				disabled={isDownloading}
			>
				<ArrowLeftIcon class="mr-2 h-4 w-4" />
				Back
			</Button>
			<div class="flex flex-wrap gap-2">
				<button
					id="status"
					class="flex min-w-[60px] cursor-pointer items-center justify-center hover:opacity-90"
					onclick={onStatusBadgeClick}
					disabled={isDownloading || isStatusUpdating}
				>
					{#if isStatusUpdating}
						<Spinner size={16} />
					{:else}
						<Badge variant={getStatusVariant(invoice.status)} class="px-3 py-1 text-sm">
							{getStatusLabel(invoice.status)}
						</Badge>
					{/if}
				</button>

				<Button
					id="edit"
					variant="default"
					size="sm"
					onclick={editInvoice}
					disabled={isDownloading || isStatusUpdating}
				>
					<EditIcon class="h-4 w-4 sm:mr-1" />
					<span class="hidden sm:inline">Edit</span>
				</Button>
				<Button
					id="download"
					variant="outline"
					size="sm"
					onclick={downloadInvoice}
					disabled={isDownloading || isStatusUpdating}
					class="justify-center"
				>
					{#if isDownloading}
						<Spinner class="h-4 w-4" size={16} />
						<span class="hidden sm:inline">Downloading...</span>
					{:else}
						<DownloadIcon class="h-4 w-4 sm:mr-1" />
						<span class="hidden sm:inline">Download</span>
					{/if}
				</Button>
				<EmailDialog
					templateData={formattedEmail}
					{onSendEmail}
					{isSending}
					authMethod={$authMethod}
				>
					<Button
						id="send"
						variant="outline"
						size="sm"
						disabled={isDownloading || isStatusUpdating}
					>
						<SendIcon class="h-4 w-4 sm:mr-1" />
						<span class="hidden sm:inline">Send Invoice</span>
					</Button>
				</EmailDialog>
			</div>
		</div>
	{/if}
</div>
