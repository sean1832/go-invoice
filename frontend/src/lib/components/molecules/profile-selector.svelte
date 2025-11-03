<script lang="ts">
	import * as Item from '$lib/components/ui/item';
	import * as Select from '$lib/components/ui/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import SettingsIcon from '@lucide/svelte/icons/settings';
	import UserPlusIcon from '@lucide/svelte/icons/user-plus';
	import ChevronDownIcon from '@lucide/svelte/icons/chevron-down';
	import ChevronUpIcon from '@lucide/svelte/icons/chevron-up';
	import type { Party } from '@/types/invoice';
	import { cn } from '$lib/utils';

	interface Profile {
		id: string;
		name: string;
		email?: string;
		abn?: string;
		address?: string;
		phone?: string;
		taxRate?: number;
		paymentInfo?: {
			method: string;
			accountName: string;
			bsb: string;
			accountNumber: string;
		};
	}

	interface Props {
		class?: string;
		type: 'provider' | 'client';
		profiles: Profile[];
		selectedProfileId?: string;
		profileData: Party;
		onSelect?: (profileId: string) => void;
		onConfigure?: () => void;
		onAddNew?: () => void;
	}

	let {
		class: customClass = '',
		type,
		profiles,
		selectedProfileId = $bindable(),
		profileData,
		onSelect,
		onConfigure,
		onAddNew
	}: Props = $props();

	// Title and description based on type
	const title = type === 'provider' ? 'Provider' : 'Client';
	const placeholder = type === 'provider' ? 'Select provider' : 'Select client';

	// Collapsed state
	let isExpanded = $state(false);

	// Track select open state
	let selectOpen = $state(false);

	function handleSelect(profileId: string) {
		selectedProfileId = profileId;

		onSelect?.(profileId);
	}

	function handleConfigure() {
		onConfigure?.();
	}

	function handleAddNew() {
		onAddNew?.();
	}

	function toggleExpanded() {
		isExpanded = !isExpanded;
	}
</script>

<Item.Root
	variant="muted"
	class={cn('h-full w-full flex-col! flex-nowrap! items-start', customClass)}
>
	<Item.Header class="w-full">
		<Item.Content class="min-w-0 shrink">
			<Item.Title>{title}</Item.Title>
		</Item.Content>
		<Item.Actions class="shrink-0">
			<div class="flex flex-nowrap items-center gap-2">
				<Select.Root
					type="single"
					bind:open={selectOpen}
					value={selectedProfileId}
					onValueChange={(v) => v && handleSelect(v)}
				>
					<Select.Trigger class="w-[200px]">
						<span>{profiles.find((p) => p.id === selectedProfileId)?.name || placeholder}</span>
					</Select.Trigger>
					<Select.Content>
						{#each profiles as profile}
							<Select.Item value={profile.id} label={profile.name}>
								{profile.name}
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
				<Button variant="outline" size="sm" onclick={handleConfigure} disabled={!selectedProfileId}>
					<SettingsIcon class="h-4 w-4" />
				</Button>
				<Button variant="outline" size="sm" onclick={handleAddNew}>
					<UserPlusIcon class="h-4 w-4" />
				</Button>
				{#if selectedProfileId}
					<Button variant="ghost" size="sm" onclick={toggleExpanded}>
						{#if isExpanded}
							<ChevronUpIcon class="h-4 w-4" />
						{:else}
							<ChevronDownIcon class="h-4 w-4" />
						{/if}
					</Button>
				{/if}
			</div>
		</Item.Actions>
	</Item.Header>
	{#if selectedProfileId}
		{#if isExpanded}
			<Item.Content class="w-full pt-2">
				<div class="grid grid-cols-2 gap-x-6 gap-y-3 text-sm wrap-break-word">
					<div>
						<span class="font-medium text-muted-foreground">Name:</span>
						<span class="ml-2">{profileData.name || 'Not configured'}</span>
					</div>
					<div>
						<span class="font-medium text-muted-foreground">ABN:</span>
						<span class="ml-2">{profileData.abn || 'Not configured'}</span>
					</div>
					<div class="col-span-2">
						<span class="font-medium text-muted-foreground">Address:</span>
						<span class="ml-2">{profileData.address || 'Not configured'}</span>
					</div>
					<div>
						<span class="font-medium text-muted-foreground">Email:</span>
						<span class="ml-2">{profileData.email || 'Not configured'}</span>
					</div>
					<div>
						<span class="font-medium text-muted-foreground">Phone:</span>
						<span class="ml-2">{profileData.phone || 'Not configured'}</span>
					</div>

					{#if type === 'client'}
						{@const selectedProfile = profiles.find((p) => p.id === selectedProfileId)}
						{#if selectedProfile?.taxRate !== undefined}
							<div class="col-span-2 mt-2 border-t pt-3">
								<div
									class="mb-2 text-xs font-semibold tracking-wide text-muted-foreground uppercase"
								>
									Tax Information
								</div>
								<div>
									<span class="font-medium text-muted-foreground">Tax Rate:</span>
									<span class="ml-2">{selectedProfile.taxRate}%</span>
								</div>
							</div>
						{/if}
					{/if}

					{#if type === 'provider'}
						{@const selectedProfile = profiles.find((p) => p.id === selectedProfileId)}
						{#if selectedProfile?.paymentInfo}
							<div class="col-span-2 mt-2 border-t pt-3">
								<div
									class="mb-2 text-xs font-semibold tracking-wide text-muted-foreground uppercase"
								>
									Payment Information
								</div>
								<div class="grid grid-cols-2 gap-x-6 gap-y-2">
									<div>
										<span class="font-medium text-muted-foreground">Method:</span>
										<span class="ml-2">{selectedProfile.paymentInfo.method}</span>
									</div>
									<div>
										<span class="font-medium text-muted-foreground">Account Name:</span>
										<span class="ml-2">{selectedProfile.paymentInfo.accountName}</span>
									</div>
									<div>
										<span class="font-medium text-muted-foreground">BSB:</span>
										<span class="ml-2">{selectedProfile.paymentInfo.bsb}</span>
									</div>
									<div>
										<span class="font-medium text-muted-foreground">Account #:</span>
										<span class="ml-2">{selectedProfile.paymentInfo.accountNumber}</span>
									</div>
								</div>
							</div>
						{/if}
					{/if}
				</div>
			</Item.Content>
		{/if}
	{:else}
		<Item.Content class="pt-2">
			<p class="text-sm text-muted-foreground italic">
				No {type} selected. Please select a {type} profile.
			</p>
		</Item.Content>
	{/if}
</Item.Root>
