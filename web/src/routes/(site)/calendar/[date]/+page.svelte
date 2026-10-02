<script lang="ts">
	import type { PageProps } from './$types';
	import { Button } from '$lib/components/ui/button';
	import { DateFormatter } from '@internationalized/date';
	import type { TimeEntry } from '$lib/types';
	import { ArrowLeft } from '@lucide/svelte';
	import ProgressCard from './ProgressCard.svelte';
	import RegisterCard from './RegisterCard.svelte';
	import TimeEntryOverview from './TimeEntryOverview.svelte';
	import HourPercentSwitch from './HourPercentSwitch.svelte';
	import { localStore } from '$lib/stores';
	import { Hour } from '$lib/utils';
	import { format, startOfMonth } from 'date-fns';
	import * as Navbar from '$lib/components/navbar';

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

	let usePercentStore = localStore('usePercent', false);

	// the live timer tracks time right now, so it is only offered for today (or while a timer is active)
	const showLiveTimer = $derived(
		data.timer != null || data.createForm.data.date == format(new Date(), 'yyyy-MM-dd')
	);
</script>

<Navbar.Root>
	<Navbar.Action>
		<Button
			href={`/calendar?date=${format(startOfMonth(daySummary.date), 'yyyy-MM-dd')}`}
			variant="link"
		>
			<ArrowLeft />
			Tilbage
		</Button>
	</Navbar.Action>
	<Navbar.Title>{formattedDate}</Navbar.Title>
	<Navbar.Action side="right">
		<HourPercentSwitch
			bind:value={$usePercentStore}
			onActiveChange={(active) => usePercentStore.set(active == 'right')}
		/>
	</Navbar.Action>
</Navbar.Root>

<div class="grid w-full grid-cols-2 gap-6">
	<div class="flex flex-col md:flex-row-reverse col-span-2 gap-6 [&>*]:w-full">
		<RegisterCard
			maxHours={daySummary.maxHours / Hour}
			formData={data.createForm}
			categories={data.categories}
			usePercent={$usePercentStore}
			timer={data.timer}
			showLive={showLiveTimer}
		/>
		<ProgressCard {daySummary} usePercent={$usePercentStore} />
	</div>
	<TimeEntryOverview {daySummary} {formattedDate} usePercent={$usePercentStore} />
</div>
