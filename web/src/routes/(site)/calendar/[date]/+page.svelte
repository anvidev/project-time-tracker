<script lang="ts">
	import type { PageProps } from './$types';
	import { Button } from '$lib/components/ui/button';
	import { DateFormatter } from '@internationalized/date';
	import type { TimeEntry } from '$lib/types';
	import { Hour } from '$lib/utils';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { ArrowLeft } from "@lucide/svelte";
	import ProgressCard from './ProgressCard.svelte';
	import TimeEntryForm from './TimeEntryForm.svelte';

	let { data }: PageProps = $props();

	const daySummary = $derived.by(() => {
		data.daySummary.timeEntries.forEach((entry) => console.log(entry.duration));
		const grouped: Record<number, TimeEntry> = {};
		for (const entry of data.daySummary.timeEntries) {
			if (grouped[entry.categoryId] == undefined) {
				grouped[entry.categoryId] = {
					...entry
				};
			} else {
				grouped[entry.categoryId].duration += entry.duration;
			}
		}

		return {
			...data.daySummary,
			timeEntries: Array.from(Object.values(grouped)).sort((a, b) => b.duration - a.duration)
		};
	});

	const dateFormatter = new DateFormatter('da-DK', { dateStyle: 'long' });

	const formattedDate = $derived(
		dateFormatter.format(
			new Date(daySummary.date)
		)
	)
</script>

<div class="flex w-full items-center p-4 border rounded-xl mb-6 gap-4">
	<Button href="/calendar">
		<ArrowLeft />
		Tilbage
	</Button>
	<p class="font-semibold text-xl tracking-tight">Tidsregistrering d. {formattedDate}</p>
</div>

<div class="grid w-full grid-cols-2 gap-6">
	<ProgressCard {daySummary} />

	<TimeEntryForm formData={data.form} categories={data.categories} />

	<Card class="col-span-2">
		<CardHeader>
			<CardTitle>
				Registreringer d. {formattedDate}
			</CardTitle>
			<CardDescription
				>Timer: {daySummary.totalHours / Hour} Max timer: {daySummary.maxHours /
					Hour}</CardDescription
			>
		</CardHeader>
		<CardContent>
			{#each daySummary.timeEntries as entry (entry.id)}
				<div class="bg-muted mb-2 flex items-center justify-between rounded-lg border p-2">
					<div>
						<p class="font-semibold tracking-tight">{entry.category}</p>
						<p class="text-muted-foreground text-sm">
							{((entry.duration / daySummary.totalHours) * 100).toFixed(1)}% af total
						</p>
					</div>
					<div>
						<Badge class="bg-background" variant="outline">
							{(entry.duration / Hour).toFixed(2)}t
						</Badge>
					</div>
				</div>
			{/each}
		</CardContent>
	</Card>
</div>
