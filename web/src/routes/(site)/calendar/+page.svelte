<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import type { PageProps } from './$types';
	import type { WeekDay } from '$lib/types';
	import CalendarDayLink from './CalendarDayLink.svelte';
	import MonthPicker from './MonthPicker.svelte';
	import { User } from '@lucide/svelte';
	import { buttonVariants } from '$lib/components/ui/button';

	const { data }: PageProps = $props();

	const { summary, calendar, date } = $derived(data);

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

		const remainingDays = 7 - ((days.length + weekDaysBeforeFirst) % 7);

		return [
			...Array.from(Array(weekDaysBeforeFirst)),
			...days,
			...Array.from(Array(remainingDays))
		];
	});
</script>

<Card.Root class="my-6 max-h-[90dvh] w-full">
	<Card.Header>
		<MonthPicker {date} />
		<Card.Action>
			<a
				href="/me"
				class={buttonVariants({ variant: 'outline', size: 'icon', class: 'cursor-pointer' })}
			>
				<User />
			</a>
		</Card.Action>
	</Card.Header>
	<Card.Content class="grid w-full grid-cols-7 gap-2">
		{#each Object.keys(weekDayMap) as weekDay}
			<p class="text-muted-foreground w-full text-center text-sm font-semibold">{weekDay}</p>
		{/each}
		{#each daysWithCalendarInfo as day (day?.id ?? crypto.randomUUID())}
			<CalendarDayLink {day} />
		{/each}
	</Card.Content>
</Card.Root>
