<script lang="ts">
	import { cn, toDurationString } from '$lib/utils';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Eye } from '@lucide/svelte';
	import type { TimeEntry } from '$lib/types';

	const { entries }: { entries: TimeEntry[] } = $props();
</script>

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
		{#each entries as entry (entry)}
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
				<Table.Cell class="w-20">
					<Eye class="ml-auto size-4" />
				</Table.Cell>
			</Table.Row>
		{/each}
	</Table.Body>
</Table.Root>
