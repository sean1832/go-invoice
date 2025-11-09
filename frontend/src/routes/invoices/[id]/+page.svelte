<script lang="ts">
	import { Badge } from '@/components/ui/badge';
	import Button from '@/components/ui/button/button.svelte';
	import InvoiceDisplayCard from '@/components/organisms/invoice-display/invoice-display-card.svelte';
	import type { EmailConfig, EmailTemplate, Invoice } from '@/types/invoice';
	import EditIcon from '@lucide/svelte/icons/pencil';
	import DownloadIcon from '@lucide/svelte/icons/download';
	import SendIcon from '@lucide/svelte/icons/send';
	import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';
	import ErrorAlert from '@/components/molecules/error-alert.svelte';
	import { api } from '@/services';
	import EmailDialog from '@/components/molecules/email-dialog.svelte';
	import { formatEmailTemplate, validateEmailConfig } from '$lib/helpers';

	interface Props {
		data: {
			invoice: Invoice | null;
			emailConfig: EmailTemplate | null;
			error?: string;
		};
	}

	let { data }: Props = $props();
	let invoice = $derived(data.invoice as Invoice);
	let emailData = $derived(data.emailConfig as EmailTemplate);
	let error = $derived(data.error);

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

	async function downloadInvoice() {
		try {
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
		}
	}

	let isSending = $state(false);

	async function onSubmit(emailConfig: EmailConfig) {
		// validate first
		const validation = validateEmailConfig(emailConfig);
		if (!validation.isValid) {
			alert(
				'Email validation failed, please fix the following errors:\n\n' +
					validation.errors.join('\n')
			);
			return;
		}

		try {
			isSending = true;
			await api.invoices.sendInvoiceEmail(fetch, invoice.id, emailConfig);
			// Show success message
			alert('Invoice email sent successfully.');
		} catch (error) {
			console.error('Error sending invoice email:', error);
			alert(
				error instanceof Error ? error.message : 'Failed to send invoice email. Please try again.'
			);
		} finally {
			isSending = false;
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
		<InvoiceDisplayCard {invoice} />

		<!-- Action Buttons Bar -->
		<div class="mt-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
			<Button variant="ghost" size="sm" onclick={goBack} class="self-start">
				<ArrowLeftIcon class="mr-2 h-4 w-4" />
				Back
			</Button>
			<div class="flex flex-wrap gap-2">
				<Badge variant={getStatusVariant(invoice.status)} class="px-3 py-1 text-sm">
					{getStatusLabel(invoice.status)}
				</Badge>
				<Button variant="default" size="sm" onclick={editInvoice}>
					<EditIcon class="mr-2 h-4 w-4" />
					<span class="hidden sm:inline">Edit</span>
				</Button>
				<Button variant="outline" size="sm" onclick={downloadInvoice}>
					<DownloadIcon class="h-4 w-4 sm:mr-2" />
					<span class="hidden sm:inline">Download PDF</span>
				</Button>
				<EmailDialog templateData={formattedEmail} {onSubmit} {isSending}>
					<Button variant="outline" size="sm">
						<SendIcon class="mr-2 h-4 w-4" />
						<span class="hidden sm:inline">Send Invoice</span>
					</Button>
				</EmailDialog>
			</div>
		</div>
	{/if}
</div>
