<script lang="ts">
	import type { SummaryDay } from '$lib/types';
	import { Hour, maxFractionDigits } from '$lib/utils';
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

	const { timeEntries, totalHours: totalHoursProp, maxHours: maxHoursProp } = $derived(daySummary);

	const totalHours = $derived.by(() => {
		if (totalHoursProp == 0) {
			return 0;
		}

		return maxFractionDigits(totalHoursProp / Hour, 2);
	});

	const maxHours = $derived.by(() => {
		if (maxHoursProp == 0) {
			return 0;
		}

		return maxFractionDigits(maxHoursProp / Hour, 2);
	});
</script>

{#if timeEntries.length > 0}
	<Card class="col-span-2">
		<CardHeader>
			<CardTitle>
				Registreringer d. {formattedDate}
			</CardTitle>
			<CardDescription>
				{totalHours} af {maxHours} timer registreret
			</CardDescription>
		</CardHeader>
		<CardContent>
			{#each timeEntries as entry (entry.id)}
				<TimeEntryCard {entry} {usePercent} maxHours={daySummary.maxHours} />
			{/each}
		</CardContent>
	</Card>
{/if}
