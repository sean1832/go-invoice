<script lang="ts">
  import UserIcon from  "@lucide/svelte/icons/circle-user-round"
	import Button from "@/components/ui/button/button.svelte";

</script>

<nav class="p-4">
  <div class="container mx-auto flex justify-between">
    <h1 class="text-2xl font-bold">GoInvoice</h1>
    <div class="space-x-4">
      <Button variant="ghost" class="cursor-pointer">Invoices</Button>
      <Button variant="ghost" class="cursor-pointer">Clients</Button>
      <Button variant="ghost" class="cursor-pointer">Settings</Button>
      <Button class="p-1 rounded-full cursor-pointer" variant="outline">
        <UserIcon class="w-6 h-6" />
      </Button>
    </div>
   
  </div>
  <!-- Divider -->
  <div class="border-t border-gray-800 my-4"></div>
</nav>