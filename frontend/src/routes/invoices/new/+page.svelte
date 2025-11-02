<script lang="ts">
	import InvoiceForm from '@/components/custom/invoice-form.svelte';

	function handleSave(data: any) {
		// TODO: Implement API call to save invoice
		console.log('Saving new invoice:', data);
		// Redirect to invoice list after successful save
		window.location.href = '/';
	}

	function handleCancel() {
		window.location.href = '/';
	}
</script>

<div class="container mx-auto p-4 max-w-5xl">
	<InvoiceForm mode="create" onSave={handleSave} onCancel={handleCancel} />
</div>
