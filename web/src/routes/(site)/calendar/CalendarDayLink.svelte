<script lang="ts">
	import { maxFractionDigits } from "$lib/utils";
	import type { SummaryDay, WeekDay } from '$lib/types';
	import { parseDate } from '@internationalized/date';
	import { isFuture, isPast } from 'date-fns';
	interface ExtendedSummaryDay extends SummaryDay {
		weekday: WeekDay;
		holliday: boolean;
		dayName: string;
		isWeekend: boolean;
	}

	const { day }: { day: ExtendedSummaryDay | undefined } = $props();

	const progress = $derived.by(() =>
		day != undefined ? Math.min((day.totalHours / day.maxHours) * 100, 100) : 0
	);

	const getProgressColor = () => {
		if (progress < 25) return 'stroke-red-500';
		if (progress < 50) return 'stroke-orange-500';
		if (progress < 75) return 'stroke-yellow-500';
		if (progress < 100) return 'stroke-green-300';
		return 'stroke-green-500';
	};

	const getProgressBgColor = () => {
		return !isFuture(day?.date ?? new Date()) && progress < 1 ? 'text-red-500' : 'text-slate-200';
	};

	const commonStyles = 'flex w-full aspect-16/8 flex-col justify-between rounded-xl border p-2';
</script>

{#if day == undefined}
	<div class={`bg-muted ${commonStyles}`}></div>
{:else if day.holliday}
	<div class={`bg-muted-foreground/30 bg-striped ${commonStyles}`}>
		<p class="text-muted-foreground w-full space-x-2 text-sm font-semibold">
			<span>{parseDate(day.date).day}</span>
			<span>{day.dayName}</span>
		</p>
	</div>
{:else if day.isWeekend}
	<div class={`bg-striped bg-muted/25 ${commonStyles}`}>
		<p class="text-muted-foreground w-full space-x-2 text-sm font-semibold">
			<span>{parseDate(day.date).day}</span>
		</p>
	</div>
{:else}
	<a
		href={`/calendar/${day.date}`}
		class={`${commonStyles} transition-all hover:border-blue-800 hover:shadow-sm`}
	>
		<div class="flex justify-between">
			<p class="text-muted-foreground w-full text-sm font-semibold">
				{parseDate(day.date).day}
			</p>
			<div class="flex gap-1 items-end justify-center">
				{#if isPast(day.date)}
					<p class="text-muted-foreground text-center text-xs">
						{maxFractionDigits(progress, 2)}%
					</p>
				{/if}
				<svg class="size-4 -rotate-90 transform" viewBox="0 0 100 100">
					<circle
						cx="50"
						cy="50"
						r="45"
						stroke="currentColor"
						stroke-width="8"
						fill="transparent"
						class={`${getProgressBgColor()}`}
					/>
					<circle
						cx="50"
						cy="50"
						r="45"
						stroke="currentColor"
						stroke-width="8"
						fill="transparent"
						stroke-dasharray={`${2 * Math.PI * 45}`}
						stroke-dashoffset={`${2 * Math.PI * 45 * (1 - progress / 100)}`}
						class={`${getProgressColor()} transition-all duration-1000 ease-out`}
						stroke-linecap="round"
					/>
				</svg>
			</div>
		</div>
	</a>
{/if}
