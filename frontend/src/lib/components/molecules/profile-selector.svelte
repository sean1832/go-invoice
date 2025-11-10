<script lang="ts">
	import * as Item from '$lib/components/ui/item';
	import * as Select from '$lib/components/ui/select';
	import * as Popover from '$lib/components/ui/popover';
	import Button from '$lib/components/ui/button/button.svelte';
	import SettingsIcon from '@lucide/svelte/icons/settings';
	import UserPlusIcon from '@lucide/svelte/icons/user-plus';
	import InfoIcon from '@lucide/svelte/icons/info';
	import type { Party } from '@/types/invoice';
	import { cn } from '$lib/utils';

	interface Profile {
		id: string;
		name: string;
		email?: string;
		abn?: string;
		address?: string;
		phone?: string;
		tax_rate?: number;
		payment_info?: {
			method: string;
			account_name: string;
			bsb: string;
			account_number: string;
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

	// Track select open state
	let selectOpen = $state(false);

	// Local storage key for caching
	const CACHE_KEY =
		type === 'client' ? 'invoice_last_selected_client' : 'invoice_last_selected_provider';

	// Save selection to localStorage
	function saveToCache(profileId: string) {
		if (typeof window !== 'undefined' && window.localStorage) {
			try {
				localStorage.setItem(CACHE_KEY, profileId);
			} catch (error) {
				console.warn('Failed to save profile selection to cache:', error);
			}
		}
	}

	// Load selection from localStorage
	function loadFromCache(): string | null {
		if (typeof window !== 'undefined' && window.localStorage) {
			try {
				return localStorage.getItem(CACHE_KEY);
			} catch (error) {
				console.warn('Failed to load profile selection from cache:', error);
				return null;
			}
		}
		return null;
	}

	function handleSelect(profileId: string) {
		selectedProfileId = profileId;
		saveToCache(profileId);
		onSelect?.(profileId);
	}

	function handleConfigure() {
		onConfigure?.();
	}

	function handleAddNew() {
		onAddNew?.();
	}

	// Export function for parent components to load cached selection
	export function getCachedSelection(): string | null {
		return loadFromCache();
	}
</script>

<Item.Root
	variant="muted"
	class={cn('h-full w-full flex-col! flex-nowrap! items-start', customClass)}
>
	<Item.Header class="w-full flex-col! items-start! gap-3">
		<Item.Content class="w-full">
			<Item.Title>{title}</Item.Title>
		</Item.Content>
		<Item.Actions class="w-full">
			<div class="flex w-full flex-nowrap items-center gap-2">
				<Select.Root
					type="single"
					bind:open={selectOpen}
					value={selectedProfileId}
					onValueChange={(v) => v && handleSelect(v)}
				>
					<Select.Trigger class="w-[280px] md:w-[200px] lg:w-[330px]">
						<span class="block truncate">
							{profiles.find((p) => p.id === selectedProfileId)?.name || placeholder}
						</span>
					</Select.Trigger>
					<Select.Content>
						{#each profiles as profile}
							<Select.Item value={profile.id} label={profile.name}>
								<span class="block truncate" title={profile.name}>{profile.name}</span>
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
					<Popover.Root>
						<Popover.Trigger
							class="inline-flex h-8 items-center justify-center gap-2 rounded-md px-3 text-sm font-medium whitespace-nowrap transition-colors hover:bg-accent hover:text-accent-foreground focus-visible:ring-1 focus-visible:ring-ring focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50"
						>
							<InfoIcon class="h-4 w-4" />
						</Popover.Trigger>
						<Popover.Content class="w-[400px]">
							<div class="grid gap-4">
								<div class="space-y-2">
									<h4 class="leading-none font-medium">{title} Details</h4>
								</div>
								<div class="grid gap-3 text-sm">
									<div class="grid grid-cols-[100px_1fr] items-start gap-2">
										<span class="font-medium text-muted-foreground">Name:</span>
										<span class="wrap-break-word">{profileData.name || 'Not configured'}</span>
									</div>
									<div class="grid grid-cols-[100px_1fr] items-start gap-2">
										<span class="font-medium text-muted-foreground">ABN:</span>
										<span>{profileData.abn || 'Not configured'}</span>
									</div>
									<div class="grid grid-cols-[100px_1fr] items-start gap-2">
										<span class="font-medium text-muted-foreground">Address:</span>
										<span class="wrap-break-word">{profileData.address || 'Not configured'}</span>
									</div>
									<div class="grid grid-cols-[100px_1fr] items-start gap-2">
										<span class="font-medium text-muted-foreground">Email:</span>
										<span class="wrap-break-word">{profileData.email || 'Not configured'}</span>
									</div>
									<div class="grid grid-cols-[100px_1fr] items-start gap-2">
										<span class="font-medium text-muted-foreground">Phone:</span>
										<span>{profileData.phone || 'Not configured'}</span>
									</div>

									{#if type === 'client'}
										{@const selectedProfile = profiles.find((p) => p.id === selectedProfileId)}
										{#if selectedProfile?.tax_rate !== undefined}
											<div class="border-t pt-3">
												<div
													class="mb-2 text-xs font-semibold tracking-wide text-muted-foreground uppercase"
												>
													Tax Information
												</div>
												<div class="grid grid-cols-[100px_1fr] items-start gap-2">
													<span class="font-medium text-muted-foreground">Tax Rate:</span>
													<span>{selectedProfile.tax_rate}%</span>
												</div>
											</div>
										{/if}
									{/if}

									{#if type === 'provider'}
										{@const selectedProfile = profiles.find((p) => p.id === selectedProfileId)}
										{#if selectedProfile?.payment_info}
											<div class="border-t pt-3">
												<div
													class="mb-2 text-xs font-semibold tracking-wide text-muted-foreground uppercase"
												>
													Payment Information
												</div>
												<div class="grid gap-2">
													<div class="grid grid-cols-[100px_1fr] items-start gap-2">
														<span class="font-medium text-muted-foreground">Method:</span>
														<span>{selectedProfile.payment_info.method}</span>
													</div>
													<div class="grid grid-cols-[100px_1fr] items-start gap-2">
														<span class="font-medium text-muted-foreground">Account Name:</span>
														<span class="wrap-break-word"
															>{selectedProfile.payment_info.account_name}</span
														>
													</div>
													<div class="grid grid-cols-[100px_1fr] items-start gap-2">
														<span class="font-medium text-muted-foreground">BSB:</span>
														<span>{selectedProfile.payment_info.bsb}</span>
													</div>
													<div class="grid grid-cols-[100px_1fr] items-start gap-2">
														<span class="font-medium text-muted-foreground">Account #:</span>
														<span>{selectedProfile.payment_info.account_number}</span>
													</div>
												</div>
											</div>
										{/if}
									{/if}
								</div>
							</div>
						</Popover.Content>
					</Popover.Root>
				{/if}
			</div>
		</Item.Actions>
	</Item.Header>
</Item.Root>
