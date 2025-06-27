<script lang="ts">
	import type { SummaryDay } from '$lib/types';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Hour, maxFractionDigits, toDurationString } from '$lib/utils';
	import ProgressContent from '$lib/components/progress-content.svelte';

	const { daySummary, usePercent }: { daySummary: SummaryDay; usePercent: boolean } = $props();

	const averageTime = $derived.by(() => {
		if (daySummary.totalHours == 0 || daySummary.timeEntries.length == 0) {
			return 0;
		}

		return maxFractionDigits(daySummary.totalHours / daySummary.timeEntries.length, 2);
	});

	const completedPercentage = $derived(
		daySummary.totalHours < 0.1 || daySummary.maxHours < 0.1
			? 0
			: maxFractionDigits((daySummary.totalHours / daySummary.maxHours) * 100, 2)
	);

	const averagePercentage = $derived.by(() => {
		if (daySummary.totalHours == 0 || daySummary.timeEntries.length == 0) {
			return 0;
		}

		return maxFractionDigits(completedPercentage / daySummary.timeEntries.length, 2);
	});

	const totalHours = $derived(
		daySummary.totalHours > 0.0 ? maxFractionDigits(daySummary.totalHours / Hour, 1) : 0
	);
	const maxHours = $derived(
		daySummary.maxHours > 0.0 ? maxFractionDigits(daySummary.maxHours / Hour, 1) : 0
	);

	const totalHoursStr = $derived(toDurationString(daySummary.totalHours));

	const maxHoursStr = $derived(toDurationString(daySummary.maxHours));

	const remainingHours = $derived(
		daySummary.totalHours < 0.1 || daySummary.maxHours < 0.1
			? 0
			: Math.max(daySummary.maxHours - daySummary.totalHours, 0)
	);

	const progressPercentage = $derived(
		totalHours < 0.1 || maxHours < 0.1
			? 0
			: maxFractionDigits(Math.min((totalHours / maxHours) * 100, 100), 2)
	);
</script>

<Card>
	<CardHeader class="text-center">
		<CardTitle class="text-2xl">Dagens opsummering</CardTitle>
	</CardHeader>
	<CardContent class="flex h-full flex-col items-center gap-8">
		<ProgressContent
			{completedPercentage}
			{progressPercentage}
			{usePercent}
			{averageTime}
			{averagePercentage}
			{remainingHours}
			{totalHoursStr}
			{maxHoursStr}
			entryCount={daySummary.timeEntries.length}
		/>
	</CardContent>
</Card>
