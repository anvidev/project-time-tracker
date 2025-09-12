<script lang="ts">
	import { toDurationString } from '$lib/utils';
	import * as Table from '$lib/components/ui/table/index.js';
	import { ArrowLeft, ArrowRight } from '@lucide/svelte';
	import type { TimeEntry } from '$lib/types';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Button } from '$lib/components/ui/button/index.js';

	const sizes = ['10', '25', '50'];

	const { entries }: { entries: TimeEntry[] } = $props();

	const data = $derived(
		entries.reduce((acc, cur) => {
			const category = cur.category;
			const duration = cur.duration;

			if (!acc.has(category)) {
				acc.set(category, {
					duration: 0,
					occurence: 0
				});
			}

			const categoryEntry = acc.get(category);
			const newData = {
				duration: (categoryEntry?.duration || 0) + duration,
				occurence: (categoryEntry?.occurence || 0) + 1
			};
			acc.set(category, newData);

			return acc;
		}, new Map<string, { duration: number; occurence: number }>())
	);

	const sorted = $derived(
		Array.from(data)
			.map(([key, value]) => ({
				category: key,
				occurrence: value.occurence,
				duration: value.duration
			}))
			.sort((a, b) => b.duration - a.duration)
	);

	let pageSize = $state(sizes[0]);
	let page = $state(1);
	let maxPages = $derived(Math.floor(sorted.length / Number(pageSize) + 1));
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
				<Table.Head class="">Kategori</Table.Head>
				<Table.Head class="w-44">Antal registreringer</Table.Head>
				<Table.Head class="w-44 text-right">Timer</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each sorted.slice((page - 1) * Number(pageSize), page * Number(pageSize)) as row (row.category)}
				<Table.Row>
					<Table.Cell class="">{row.category}</Table.Cell>
					<Table.Cell class="w-44">{row.occurrence}</Table.Cell>
					<Table.Cell class="w-44 text-right">{toDurationString(row.duration)}</Table.Cell>
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
		<Button disabled={page <= 1} variant="outline" size="icon" onclick={() => navigate(-1)}>
			<ArrowLeft />
		</Button>
		<p class="text-muted-foreground text-sm">Side {page} af {maxPages}</p>
		<Button disabled={page >= maxPages} variant="outline" size="icon" onclick={() => navigate(1)}>
			<ArrowRight />
		</Button>
	</div>
</div>
