<script lang="ts">
	import * as Item from '$lib/components/ui/item';
	import * as Select from '$lib/components/ui/select';
	import Button from '$lib/components/ui/button/button.svelte';
	import SettingsIcon from '@lucide/svelte/icons/settings';
	import UserPlusIcon from '@lucide/svelte/icons/user-plus';
	import type { Party } from '@/types/invoice';

	interface Profile {
		id: string;
		name: string;
		email: string;
		abn: string;
		address?: string;
		phone?: string;
	}

	interface Props {
		type: 'provider' | 'client';
		profiles: Profile[];
		selectedProfileId?: string;
		profileData: Party;
		onSelect?: (profileId: string) => void;
		onConfigure?: () => void;
		onAddNew?: () => void;
	}

	let {
		type,
		profiles,
		selectedProfileId = $bindable(),
		profileData,
		onSelect,
		onConfigure,
		onAddNew
	}: Props = $props();

	// Title and description based on type
	const title = type === 'provider' ? 'Provider (From)' : 'Client (Bill To)';
	const description =
		type === 'provider'
			? 'Select your business profile for this invoice'
			: 'Select a client profile for this invoice';
	const placeholder = type === 'provider' ? 'Select provider' : 'Select client';

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
</script>

<Item.Root variant="muted">
	<Item.Header>
		<Item.Content>
			<Item.Title>{title}</Item.Title>
			<Item.Description>{description}</Item.Description>
		</Item.Content>
		<Item.Actions>
			<div class="flex gap-2 items-center">
				<Select.Root type="single">
					<Select.Trigger class="w-[200px]">
						<span
							>{profiles.find((p) => p.id === selectedProfileId)?.name || placeholder}</span
						>
					</Select.Trigger>
					<Select.Content>
						{#each profiles as profile}
							<Select.Item
								value={profile.id}
								label={profile.name}
								onclick={() => handleSelect(profile.id)}
							>
								{profile.name}
							</Select.Item>
						{/each}
					</Select.Content>
				</Select.Root>
				<Button variant="outline" size="sm" onclick={handleConfigure}>
					<SettingsIcon class="h-4 w-4" />
				</Button>
				<Button variant="outline" size="sm" onclick={handleAddNew}>
					<UserPlusIcon class="h-4 w-4" />
				</Button>
			</div>
		</Item.Actions>
	</Item.Header>
	{#if selectedProfileId}
		<Item.Content class="pt-2">
			<div class="grid grid-cols-2 gap-x-6 gap-y-3 text-sm">
				<div>
					<span class="text-muted-foreground font-medium">Name:</span>
					<span class="ml-2">{profileData.name || 'Not configured'}</span>
				</div>
				<div>
					<span class="text-muted-foreground font-medium">ABN:</span>
					<span class="ml-2">{profileData.abn || 'Not configured'}</span>
				</div>
				<div class="col-span-2">
					<span class="text-muted-foreground font-medium">Address:</span>
					<span class="ml-2">{profileData.address || 'Not configured'}</span>
				</div>
				<div>
					<span class="text-muted-foreground font-medium">Email:</span>
					<span class="ml-2">{profileData.email || 'Not configured'}</span>
				</div>
				<div>
					<span class="text-muted-foreground font-medium">Phone:</span>
					<span class="ml-2">{profileData.phone || 'Not configured'}</span>
				</div>
			</div>
		</Item.Content>
	{:else}
		<Item.Content class="pt-2">
			<p class="text-sm text-muted-foreground italic">
				No {type} selected. Please select a {type} profile.
			</p>
		</Item.Content>
	{/if}
</Item.Root>