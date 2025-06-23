<script lang="ts">
	import type { PageProps } from './$types';
	import { Button } from '$lib/components/ui/button';
	import { DateFormatter } from '@internationalized/date';
	import type { TimeEntry } from '$lib/types';
	import { ArrowLeft } from '@lucide/svelte';
	import ProgressCard from './ProgressCard.svelte';
	import TimeEntryForm from './TimeEntryForm.svelte';
	import TimeEntryOverview from './TimeEntryOverview.svelte';

	let { data }: PageProps = $props();

	const daySummary = $derived.by(() => {
		const grouped: Record<number, TimeEntry> = {};
		for (const entry of data.daySummary.timeEntries) {
			if (grouped[entry.id] == undefined) {
				grouped[entry.id] = {
					...entry
				};
			} else {
				grouped[entry.id].duration += entry.duration;
			}
		}

		return {
			...data.daySummary,
			timeEntries: Array.from(Object.values(grouped)).sort((a, b) => b.duration - a.duration)
		};
	});

	const dateFormatter = new DateFormatter('da-DK', { dateStyle: 'long' });

	const formattedDate = $derived(dateFormatter.format(new Date(daySummary.date)));
</script>

<div
	class="relative mb-6 flex h-[70px] w-full items-center gap-4 rounded-xl border p-4 text-center"
>
	<Button class="absolute top-1/2 left-4 -translate-y-1/2" href="/calendar" variant="link">
		<ArrowLeft />
		Tilbage
	</Button>
	<p class="w-full text-xl font-semibold tracking-tight">Tidsregistrering d. {formattedDate}</p>
</div>

<div class="grid w-full grid-cols-2 gap-6">
	<ProgressCard {daySummary} />

	<TimeEntryForm formData={data.createForm} categories={data.categories} />

	<TimeEntryOverview {daySummary} {formattedDate} />
</div>
