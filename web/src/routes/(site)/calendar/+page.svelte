<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { parseDate } from '@internationalized/date';
	import type { PageProps } from './$types';
	import { Hour } from '$lib/utils.js';

	const { data }: PageProps = $props();

	const { summary } = data;
</script>

<Card.Root class="my-6">
	<Card.Header>
		<Card.Title class="capitalize">{summary.month}</Card.Title>
	</Card.Header>
	<Card.Content class="grid grid-cols-5 gap-2">
		{#each summary.days as day}
			<a
				href={`/calendar/${day.date}`}
				class="hover:bg-muted flex size-32 flex-col justify-between rounded-xl border p-2 transition-all hover:shadow-sm"
			>
				<p class="w-full">{parseDate(day.date).day}</p>
				{#if day.totalHours / Hour >= 0.1}
					<p class="w-full text-center text-sm">{day.totalHours / Hour}t</p>
				{/if}
			</a>
		{/each}
	</Card.Content>
</Card.Root>
