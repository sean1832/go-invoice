<script lang="ts">
	import InvoiceForm from '@/components/custom/invoice-editor/invoice-editor.svelte';

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

<div class="container mx-auto max-w-5xl p-4">
	<InvoiceForm mode="create" onSave={handleSave} onCancel={handleCancel} />
</div>
