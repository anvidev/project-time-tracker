<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import type { TimeEntry } from '$lib/types';
	import { Hour, maxFractionDigits, toDurationString } from '$lib/utils';
	import { Check, Pencil, X } from '@lucide/svelte';
	import DeleteEntryModal from './DeleteEntryModal.svelte';
	import { Input } from '$lib/components/ui/input';
	import { superForm } from 'sveltekit-superforms';

	const {
		entry,
		maxHours,
		usePercent
	}: { entry: TimeEntry; maxHours: number; usePercent: boolean } = $props();

	let editable = $state(false);
	let percentState = $state(maxFractionDigits((entry.duration / maxHours) * 100, 1));

	//const formData = $derived.by(() => ());

	const formData: {
		id: number;
		durationHours: number | string;
		description: string;
	} = {
		id: entry.id,
		durationHours: toDurationString(entry.duration),
		description: entry.description
	};

	const { form, enhance } = superForm(formData, {
		dataType: 'json',
		onUpdated: ({ form: updatedForm }) => {
			$form.durationHours = updatedForm.data.durationHours;
			$form.description = updatedForm.data.description;

			editable = false;
		}
	});
</script>

<div class="bg-muted mb-2 flex items-center justify-between rounded-lg border p-2">
	<div>
		<p class="font-semibold tracking-tight">{entry.category}</p>
		<p class="text-muted-foreground text-sm">
			{entry.description ? entry.description : 'Ingen beskrivelse'}
		</p>
	</div>
	{#if !editable}
		<div class="grid grid-cols-[58px_58px_36px_36px] gap-2">
			{#if entry.duration > 0 && maxHours > 0}
				<Badge class="bg-background w-full px-3 py-2" variant="outline">
					{maxFractionDigits((entry.duration / maxHours) * 100, 2)}%
				</Badge>
			{:else}
				<div></div>
			{/if}
			<Badge class="bg-background w-full px-3 py-2" variant="outline">
				{toDurationString(entry.duration)}
			</Badge>
			<Button
				class="cursor-pointer"
				variant="outline"
				size="icon"
				onclick={() => (editable = true)}
			>
				<Pencil class="size-4" />
			</Button>
			<DeleteEntryModal id={entry.id} />
		</div>
	{:else}
		<form
			method="POST"
			action="?/updateTimeEntry"
			class="grid grid-cols-[58px_58px_36px_36px] gap-2"
			use:enhance
		>
			<Input type="hidden" bind:value={$form.description} />
			{#if usePercent}
				<div class="relative col-span-2 col-start-1">
					<Input
						type="number"
						step="0.1"
						inputmode="numeric"
						class="input-arrows-none w-full"
						bind:value={percentState}
						oninput={() => ($form.durationHours = (maxHours / Hour) * (percentState / 100))}
					/>
					<span
						class="text-muted-foreground absolute top-1/2 right-2 -translate-y-1/2 rounded-sm border px-1 pb-1 text-center text-xs"
					>
						%
					</span>
				</div>
			{:else}
				<div class="relative col-span-2 col-start-1">
					<Input
						type="text"
						class="input-arrows-none w-full"
						bind:value={$form.durationHours}
						pattern="(?:\d+|\d+t \d+m|0t)"
					/>
					<span
						class="text-muted-foreground absolute top-1/2 right-2 -translate-y-1/2 rounded-sm border px-1 pb-1 text-center text-xs"
					>
						t
					</span>
				</div>
			{/if}
			<Button
				type="submit"
				class="cursor-pointer text-green-600 hover:text-green-600"
				variant="outline"
				size="icon"
			>
				<Check class="size-4" />
			</Button>
			<Button
				type="button"
				class="cursor-pointer text-red-500 hover:text-red-500"
				variant="outline"
				size="icon"
				onclick={() => (editable = false)}
			>
				<X class="size-4" />
			</Button>
		</form>
	{/if}
</div>
