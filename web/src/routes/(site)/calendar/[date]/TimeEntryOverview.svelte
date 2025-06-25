<script lang="ts">
	import type { SummaryDay } from '$lib/types';
	import { Hour, maxFractionDigits, toDurationString } from '$lib/utils';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import TimeEntryCard from './TimeEntryCard.svelte';

	const {
		daySummary,
		formattedDate,
		usePercent
	}: {
		daySummary: SummaryDay;
		formattedDate: string;
		usePercent: boolean;
	} = $props();

	const { timeEntries, totalHours, maxHours } = $derived(daySummary);
</script>

{#if timeEntries.length > 0}
	<Card class="col-span-2">
		<CardHeader>
			<CardTitle>
				Registreringer d. {formattedDate}
			</CardTitle>
			<CardDescription>
				{toDurationString(totalHours)} af {toDurationString(maxHours)} registreret
			</CardDescription>
		</CardHeader>
		<CardContent>
			{#each timeEntries as entry (entry.id)}
				<TimeEntryCard {entry} {usePercent} maxHours={daySummary.maxHours} />
			{/each}
		</CardContent>
	</Card>
{/if}
