<script lang="ts">
	import { cn, toDurationString } from '$lib/utils';
	import * as Table from '$lib/components/ui/table/index.js';
	import { ArrowLeft, ArrowRight, Eye } from '@lucide/svelte';
	import type { TimeEntry } from '$lib/types';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Button, buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';

	const sizes = ['10', '25', '50', '100', '250'];

	const { entries }: { entries: TimeEntry[] } = $props();
	let pageSize = $state(sizes[0]);
	let page = $state(1);
	let maxPages = $derived(Math.floor(entries.length / Number(pageSize) + 1));
	const triggerContent = $derived(sizes.find((f) => f === pageSize) ?? 'Select a fruit');

	function navigate(amount: number) {
		const newPage = page + amount;
		if (newPage < 1) return;
		if (newPage > maxPages) return;
		page = newPage;
	}
</script>

<div class="bg-background w-full rounded-md border shadow-sm">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head class="w-28">Dato</Table.Head>
				<Table.Head class="w-44">Bruger</Table.Head>
				<Table.Head class="w-44">Kategori</Table.Head>
				<Table.Head class="hidden max-w-80 lg:table-cell">Beskrivelse</Table.Head>
				<Table.Head class="w-44 text-right">Timer</Table.Head>
				<Table.Head class="w-20 text-right"></Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each entries.slice((page - 1) * Number(pageSize), page * Number(pageSize)) as entry (entry)}
				<Table.Row>
					<Table.Cell class="w-28">{entry.date}</Table.Cell>
					<Table.Cell class="w-44">{entry.userName}</Table.Cell>
					<Table.Cell class="w-44">{entry.category}</Table.Cell>
					<Table.Cell
						class={cn(
							'hidden max-w-80 truncate lg:table-cell',
							entry.description == '' && 'text-muted-foreground'
						)}>{entry.description != '' ? entry.description : 'Ingen beskrivelse'}</Table.Cell
					>
					<Table.Cell class="w-44 text-right">{toDurationString(entry.duration)}</Table.Cell>
					<Table.Cell class="flex w-20 justify-end">
						<Dialog.Root>
							<Dialog.Trigger
								class={buttonVariants({
									variant: 'ghost',
									size: 'icon',
									class: 'size-7 items-center justify-center'
								})}
							>
								<Eye class="size-4" />
							</Dialog.Trigger>
							<Dialog.Content class="sm:max-w-[425px]">
								<Dialog.Header>
									<Dialog.Title>Beskrivelse</Dialog.Title>
								</Dialog.Header>
								<div class="grid gap-4 py-4">
									<p>{entry.description ?? 'Ingen beskrivelse'}</p>
								</div>
							</Dialog.Content>
						</Dialog.Root>
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>

<div class="flex w-full items-center justify-between">
	<div>
		<Select.Root type="single" name="favoriteFruit" bind:value={pageSize}>
			<Select.Trigger class="w-[100px]">
				{triggerContent}
			</Select.Trigger>
			<Select.Content>
				<Select.Group>
					<Select.Label>Rækker pr. side</Select.Label>
					{#each sizes as size (size)}
						<Select.Item value={size} label={size}>
							{size}
						</Select.Item>
					{/each}
				</Select.Group>
			</Select.Content>
		</Select.Root>
	</div>

	<div class="flex items-center gap-3">
		<Button variant="outline" size="icon" onclick={() => navigate(-1)}>
			<ArrowLeft />
		</Button>
		<p class="text-muted-foreground text-sm">Side {page} af {maxPages}</p>
		<Button variant="outline" size="icon" onclick={() => navigate(1)}>
			<ArrowRight />
		</Button>
	</div>
</div>
