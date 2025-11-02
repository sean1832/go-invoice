<script lang="ts">
	import UserIcon from '@lucide/svelte/icons/circle-user-round';
	import Button from '@/components/ui/button/button.svelte';
</script>

<nav class="p-4">
	<div class="container mx-auto flex justify-between">
		<h1 class="text-2xl font-bold">GoInvoice</h1>
		<div class="space-x-4">
			<Button variant="ghost" class="cursor-pointer">Invoices</Button>
			<Button variant="ghost" class="cursor-pointer">Clients</Button>
			<Button variant="ghost" class="cursor-pointer">Settings</Button>
			<Button class="cursor-pointer rounded-full p-1" variant="outline">
				<UserIcon class="h-6 w-6" />
			</Button>
		</div>
	</div>
	<!-- Divider -->
	<div class="my-4 border-t border-gray-800"></div>
</nav>
