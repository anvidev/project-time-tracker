<script lang="ts">
	import { cn, maxFractionDigits } from '$lib/utils';
	import type { SummaryDay, WeekDay } from '$lib/types';
	import { parseDate } from '@internationalized/date';
	import { isFuture, isPast, isToday } from 'date-fns';
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

	const today = $derived(day ? isToday(day.date) : false);

	const progress = $derived.by(() =>
		day != undefined ? Math.min((day.totalHours / day.maxHours) * 100, 100) : 0
	);

	const progressStr = $derived.by(() => {
		if (day == undefined || day.totalHours == 0 || day.maxHours == 0) {
			return '0%';
		}

		return `${maxFractionDigits((day.totalHours / day.maxHours) * 100, 2)}%`;
	});

	const animatedBg = new Tween(0, { delay: 100, duration: 1000, easing: cubicOut });
	const animatedProgress = new Tween(0, { delay: 100, duration: 1000, easing: cubicOut });

	const getProgressColor = () => {
		if (progress < 25) return 'fill-red-500';
		if (progress < 50) return 'fill-orange-500';
		if (progress < 75) return 'fill-yellow-500';
		if (progress < 100) return 'fill-green-300';
		return 'fill-green-500';
	};

	const getProgressTextColorForMobile = () => {
		if (progress < 25) return 'text-red-500';
		if (progress < 50) return 'text-orange-500';
		if (progress < 75) return 'text-yellow-500';
		if (progress < 100) return 'text-green-300';
		return 'text-green-500';
	}

	const getProgressDotColorForMobile = () => {
		if (progress < 25) return 'bg-red-500';
		if (progress < 50) return 'bg-orange-500';
		if (progress < 75) return 'bg-yellow-500';
		if (progress < 100) return 'bg-green-300';
		return 'bg-green-500';
	}

	const getProgressBgColor = () => {
		return 'fill-slate-200';
	};

	const commonStyles = 'flex w-full aspect-square sm:aspect-16/11 flex-col justify-between rounded-lg border py-1 px-2 md:p-2';

	onMount(() => {
		animatedProgress.target =
			progress > 0 || isFuture(day?.date ?? new Date()) || day?.holliday || day?.isWeekend
				? progress
				: 100;
		animatedBg.target = 10;
	});
</script>

{#if day == undefined}
	<div class={`bg-muted/25 border-border/50 ${commonStyles}`}></div>
{:else}
	<a
		href={`/calendar/${day.date}`}
		class={cn(`${commonStyles} relative transition-all hover:border-blue-800 hover:shadow-sm`, day.holliday && 'bg-muted/25 bg-striped',  day.isWeekend && 'bg-muted/25 bg-striped', today && 'border-blue-800 bg-blue-50/80')}
	>
		{#if !((day.holliday || day.isWeekend) && day.totalHours == 0)}
			<span class={cn("block sm:hidden absolute top-1.5 right-1.5 rounded-full size-1.5 bg-purple-500", getProgressDotColorForMobile())}></span>
		{/if}
		<div class="flex h-full justify-center items-center sm:justify-between sm:items-start">
			<div class="text-muted-foreground gap-1.5 flex items-start text-sm font-medium md:w-full">
				<div
					class={cn(
						'text-sm md:text-sm',
						today &&
							'text-blue-500 grid place-items-center tabular-nums'
					)}
				>
					{parseDate(day.date).day}
				</div>
				{#if day.holliday}
					<span class="hidden md:block text-xs md:text-sm truncate">{day.dayName}</span>
				{/if}
			</div>
			{#if !((day.holliday || day.isWeekend) && day.totalHours == 0)}
				<div class="relative h-full justify-end gap-1 hidden sm:flex">
					{#if isPast(day.date) && !(day.isWeekend || day.holliday)}
						<p class={cn('text-muted-foreground text-center text-xs md:text-sm transition-all max-md:mt-px', getProgressTextColorForMobile(), 'lg:text-foreground')}>
							{progressStr}
						</p>
					{/if}
					<div class="hidden lg:block relative w-2 h-[100%] overflow-hidden rounded-[3px]">
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
			{/if}
		</div>
	</a>
{/if}
