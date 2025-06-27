<script lang="ts">
	import * as Navbar from '$lib/components/navbar';
	import * as Card from '$lib/components/ui/card';
	import type { PageProps } from './$types';
	import type { WeekDay } from '$lib/types';
	import CalendarDayLink from './CalendarDayLink.svelte';
	import MonthPicker from './MonthPicker.svelte';
	import { ChartPie, Hourglass, User } from '@lucide/svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { toDurationString } from '$lib/utils';

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

	const [daysWithCalendarInfo, totalTime, daysWithEntries] = $derived.by(() => {
		let totalTime = 0;
		let daysWithEntries = 0;

		if (summary.days.length < 28) {
			return [[], totalTime, daysWithEntries];
		}

		const days = summary.days.map((day) => {
			const calendarDay = calendar.days.find((cDay) => day.date == cDay.date);
			if (!calendarDay) return;

			totalTime += day.totalHours;

			if (day.timeEntries.length > 0) {
				daysWithEntries++;
			}

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

		const remainingDays = 7 - ((days.length + weekDaysBeforeFirst) % 7);

		return [
			[...Array.from(Array(weekDaysBeforeFirst)), ...days, ...Array.from(Array(remainingDays))],
			totalTime,
			daysWithEntries
		];
	});

	let size = $state(0);
	let isMobile = $derived(size <= 756);
</script>

<svelte:window bind:innerWidth={size} />

<Navbar.Root>
	<Navbar.Action>
		<MonthPicker {date} />
	</Navbar.Action>
	<Navbar.Title class="hidden md:block">Kalender</Navbar.Title>
	<Navbar.Action side="right">
		<a
			href="/me"
			class={buttonVariants({ variant: 'outline', size: 'icon', class: 'ml-auto cursor-pointer' })}
		>
			<User />
		</a>
	</Navbar.Action>
</Navbar.Root>

<div class="flex w-full flex-col items-center gap-4 md:flex-row">
	<div class="flex w-full items-center gap-4 rounded-md border p-4 shadow-xs">
		<ChartPie class="mb-auto text-blue-500" />
		<div>
			<p class="text-lg leading-none font-bold">{toDurationString(totalTime)}</p>
			<small class="text-muted-foreground text-xs">Total tid denne måned</small>
		</div>
	</div>

	<div class="flex w-full items-center gap-4 rounded-md border p-4 shadow-xs">
		<Hourglass class="mb-auto text-purple-500" />
		<div>
			<p class="text-lg leading-none font-bold">
				{toDurationString(totalTime > 0 && daysWithEntries > 0 ? totalTime / daysWithEntries : 0)}
			</p>
			<small class="text-muted-foreground text-xs">Gennemsnitlig daglig tid</small>
		</div>
	</div>
</div>

<Card.Root class="w-full">
	<Card.Content class="grid w-full grid-cols-7 gap-2">
		{#each Object.keys(weekDayMap) as weekDay}
			<p class="text-muted-foreground w-full text-center text-sm font-semibold">
				{isMobile ? weekDay.substring(0, 3) : weekDay}
			</p>
		{/each}
		{#each daysWithCalendarInfo as day (day?.id ?? crypto.randomUUID())}
			<CalendarDayLink {day} />
		{/each}
	</Card.Content>
</Card.Root>
