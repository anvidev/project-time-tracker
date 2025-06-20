<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { parseDate } from '@internationalized/date';
	import type { PageProps } from './$types';
	import { Hour } from '$lib/utils.js';

	const { data }: PageProps = $props();

	const { summary, calendar } = data;

	type WeekDay = 'Mandag' | 'Tirsdag' | 'Onsdag' | 'Torsdag' | 'Fredag' | 'Lørdag' | 'Søndag';

	const weekDayMap: Record<WeekDay, number> = {
		Mandag: 0,
		Tirsdag: 1,
		Onsdag: 2,
		Torsdag: 3,
		Fredag: 4,
		Lørdag: 5,
		Søndag: 6
	};

	const daysWithCalendarInfo = $derived.by(() => {
		if (summary.days.length < 28) {
			return [];
		}

		const days = summary.days.map((day) => {
			const calendarDay = calendar.days.find((cDay) => day.date == cDay.date);
			if (!calendarDay) return;

			return {
				...day,
				weekday: calendarDay.weekday as WeekDay,
				holliday: calendarDay.holliday || calendarDay.events.some((event) => event.holliday),
				dayName:
					calendarDay.events.find((e) => Boolean(e.danishShort))?.danishShort ??
					calendarDay.dayName,
				isWeekend: weekDayMap[calendarDay.weekday as WeekDay] >= 5
			};
		});

		const first = days[0];
		const weekDaysBeforeFirst = weekDayMap[first?.weekday ?? 'Mandag'] - weekDayMap['Mandag'];

		if (weekDaysBeforeFirst == 0) {
			return days;
		}

		const remainingDays = 7-((days.length + weekDaysBeforeFirst) % 7)

		return [...Array.from(Array(weekDaysBeforeFirst)), ...days, ...Array.from(Array(remainingDays))];
	});
</script>

<Card.Root class="my-6">
	<Card.Header>
		<Card.Title class="capitalize">{summary.month}</Card.Title>
	</Card.Header>
	<Card.Content class="grid grid-cols-7 gap-2">
		{#each Object.keys(weekDayMap) as weekDay}
			<p class="text-muted-foreground w-full text-center text-sm font-semibold">{weekDay}</p>
		{/each}
		{#each daysWithCalendarInfo as day}
			{#if day == undefined}
				<div class="bg-muted flex size-28 flex-col justify-between rounded-xl border p-2"></div>
			{:else if day.holliday}
				<div class="bg-muted-foreground/30 bg-striped flex size-28 flex-col justify-between rounded-xl border p-2">
					<p class="text-muted-foreground w-full text-sm font-semibold space-x-2">
						<span>{parseDate(day.date).day}</span>
						<span>{day.dayName}</span>
					</p>
				</div>
			{:else if day.isWeekend}
				<div class="bg-striped bg-muted/25 flex size-28 flex-col justify-between rounded-xl border p-2">
					<p class="text-muted-foreground w-full text-sm font-semibold space-x-2">
						<span>{parseDate(day.date).day}</span>
					</p>
				</div>
			{:else}
				<a
					href={`/calendar/${day.date}`}
					class="flex size-28 flex-col justify-between rounded-xl border p-2 transition-all hover:border-blue-800 hover:shadow-sm"
				>
					<p class="text-muted-foreground w-full text-sm font-semibold">
						{parseDate(day.date).day}
					</p>
					{#if day.totalHours / Hour >= 0.1}
						<p class="w-full text-center text-sm">{day.totalHours / Hour}t</p>
					{/if}
				</a>
			{/if}
		{/each}
	</Card.Content>
</Card.Root>
