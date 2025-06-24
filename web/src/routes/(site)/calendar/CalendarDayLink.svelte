<script lang="ts">
	import { maxFractionDigits } from '$lib/utils';
	import type { SummaryDay, WeekDay } from '$lib/types';
	import { parseDate } from '@internationalized/date';
	import { isFuture, isPast } from 'date-fns';
	import { Tween } from 'svelte/motion';
	import { onMount } from 'svelte';
	import { cubicOut } from 'svelte/easing';
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

	const animatedBg = new Tween(0, { delay: 100, duration: 1000, easing: cubicOut });
	const animatedProgress = new Tween(0, { delay: 100, duration: 1000, easing: cubicOut });

	const getProgressColor = () => {
		if (progress < 25) return 'fill-red-500';
		if (progress < 50) return 'fill-orange-500';
		if (progress < 75) return 'fill-yellow-500';
		if (progress < 100) return 'fill-green-300';
		return 'fill-green-500';
	};

	const getProgressBgColor = () => {
		return 'fill-slate-200';
	};

	const commonStyles = 'flex w-full aspect-16/8 flex-col justify-between rounded-lg border p-2';

	onMount(() => {
		animatedProgress.target = progress > 0 || isFuture(day?.date ?? new Date()) ? progress : 100;
		animatedBg.target = 10;
	});
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
		<div class="flex h-full justify-between">
			<p class="text-muted-foreground w-full text-sm font-semibold">
				{parseDate(day.date).day}
			</p>
			<div class="relative flex h-full justify-end gap-1">
				{#if isPast(day.date)}
					<p class={`text-muted-foreground text-center text-xs transition-all`}>
						{maxFractionDigits(progress, 2)}%
					</p>
				{/if}
				<div class="relative aspect-[1/6] h-[100%] overflow-hidden rounded-[5px]">
					<svg
						class="absolute inset-0 h-full w-full"
						viewBox="0 0 10 10"
						preserveAspectRatio="none"
					>
						<rect
							x="0"
							y="0"
							width="10"
							height="10"
							fill="currentColor"
							class={`${getProgressBgColor()}`}
						/>
						<rect
							x="0"
							y={10 * (1 - animatedProgress.current / 100)}
							width="10"
							height={10 * (animatedProgress.current / 100)}
							fill="currentColor"
							class={`${getProgressColor()}`}
						/>
					</svg>
				</div>
			</div>
		</div>
	</a>
{/if}
