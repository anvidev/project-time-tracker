<script lang="ts">
	import * as Command from '$lib/components/ui/command/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import Label from '$lib/components/ui/label/label.svelte';
	import type { Category, User } from '$lib/types';
	import { Button } from '$lib/components/ui/button';
	import { CheckIcon, ChevronsUpDownIcon } from '@lucide/svelte';
	import { tick } from 'svelte';
	import { cn } from '$lib/utils';

	let { param = $bindable(), categories }: { param: number[]; categories: Category[] } = $props();

	let open = $state(false);
	let triggerRef = $state<HTMLButtonElement>(null!);

	function closeAndFocusTrigger() {
		open = false;
		tick().then(() => {
			triggerRef.focus();
		});
	}

	function getSelectedValues(selected: number[]): string {
		if (selected.length == 0) return 'Vælg kategorier';
		if (selected.length > 2) return `${selected.length} kategorier valgt`;
		return selected.map((s) => categories.find((c) => s == c.id)!.title).join(', ');
	}
</script>

<div class="grid w-full gap-1.5 lg:w-auto">
	<Label>Kategori</Label>
	<Popover.Root bind:open>
		<Popover.Trigger bind:ref={triggerRef}>
			{#snippet child({ props })}
				<Button
					variant="outline"
					{...props}
					class="justify-between"
					role="combobox"
					aria-expanded={open}
				>
					{getSelectedValues(param)}
					<ChevronsUpDownIcon class="ml-2 size-4 shrink-0 opacity-50" />
				</Button>
			{/snippet}
		</Popover.Trigger>
		<Popover.Content class="max-w-60 p-0">
			<Command.Root>
				<Command.Input placeholder="Søg Kategori..." />
				<Command.List>
					<Command.Empty>Ingen kategori fundet</Command.Empty>
					<Command.Group>
						{#each categories as category}
							<Command.Item
								value={category.id.toString()}
								onSelect={() => {
									if (param.find((p) => p == category.id)) {
										param = param.filter((p) => p != category.id);
									} else {
										param.push(category.id);
									}
									closeAndFocusTrigger();
								}}
							>
								<CheckIcon
									class={cn(
										'mr-2 size-4',
										!param.find((c) => c == category.id) && 'text-transparent'
									)}
								/>
								{category.title}
							</Command.Item>
						{/each}
					</Command.Group>
				</Command.List>
			</Command.Root>
		</Popover.Content>
	</Popover.Root>
</div>
