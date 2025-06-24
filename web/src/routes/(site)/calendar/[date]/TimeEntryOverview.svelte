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
		formattedDate
	}: {
		daySummary: SummaryDay;
		formattedDate: string;
	} = $props();

	const totalHours = $derived.by(() => {
		if (daySummary.totalHours == 0) {
			return 0;
		}

		return maxFractionDigits(daySummary.totalHours / Hour, 2);
	});

	const maxHours = $derived.by(() => {
		if (daySummary.maxHours == 0) {
			return 0;
		}

		return maxFractionDigits(daySummary.maxHours / Hour, 2);
	});
</script>

{#if daySummary.timeEntries.length > 0}
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
			{#each daySummary.timeEntries as entry (entry.id)}
				<TimeEntryCard {entry} maxHours={daySummary.maxHours} />
			{/each}
		</CardContent>
	</Card>
{/if}
