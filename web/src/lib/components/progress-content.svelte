<script lang="ts">
	import { cn, toDurationString } from '$lib/utils';
	import type { HTMLAttributes } from 'svelte/elements';

	interface Props extends Omit<HTMLAttributes<HTMLDivElement>, 'children'> {
		usePercent: boolean;
		completedPercentage: number;
		progressPercentage: number;
		averagePercentage: number;
		remainingHours: number;
		averageTime: number;
		totalHoursStr: string;
		maxHoursStr: string;
		entryCount: number;
	}

	let {
		completedPercentage,
		progressPercentage,
		usePercent,
		averageTime,
		averagePercentage,
		remainingHours,
		totalHoursStr,
		maxHoursStr,
		entryCount,
		class: className,
		...restProps
	}: Props = $props();

	const getProgressColor = () => {
		if (progressPercentage < 25) return 'stroke-red-500';
		if (progressPercentage < 50) return 'stroke-orange-500';
		if (progressPercentage < 75) return 'stroke-yellow-500';
		if (progressPercentage < 100) return 'stroke-green-300';
		return 'stroke-green-500';
	};
</script>

<div class={cn('relative h-48 w-48', className)} {...restProps}>
	<svg class="h-48 w-48 -rotate-90 transform" viewBox="0 0 100 100">
		<circle
			cx="50"
			cy="50"
			r="45"
			stroke="currentColor"
			stroke-width="8"
			fill="transparent"
			class="text-slate-200"
		/>
		<circle
			cx="50"
			cy="50"
			r="45"
			stroke="currentColor"
			stroke-width="8"
			fill="transparent"
			stroke-dasharray={`${2 * Math.PI * 45}`}
			stroke-dashoffset={`${2 * Math.PI * 45 * (1 - progressPercentage / 100)}`}
			class={`${getProgressColor()} transition-all duration-1000 ease-out`}
			stroke-linecap="round"
		/>
	</svg>

	<div class="absolute inset-0 flex flex-col items-center justify-center">
		{#if usePercent}
			<div class="text-3xl font-bold text-slate-800">
				{completedPercentage}%
			</div>
		{:else}
			<div class="text-3xl font-bold text-slate-800">
				{totalHoursStr}
			</div>
			<div class="text-sm text-slate-500">
				af {maxHoursStr}
			</div>
		{/if}
		{#if remainingHours > 0}
			<div class="mt-1 text-xs text-slate-400">
				Mangler {usePercent ? `${100 - completedPercentage}%` : toDurationString(remainingHours)}
			</div>
		{/if}
	</div>
</div>
<div class="grid w-full grid-cols-3">
	<div class="flex flex-col items-center justify-center">
		<p class="text-lg font-semibold">{entryCount}</p>
		<p class="text-muted-foreground text-sm">Registreringer</p>
	</div>
	<div class="flex flex-col items-center justify-center">
		{#if usePercent}
			<p class="text-lg font-semibold">
				{totalHoursStr}
			</p>
			<p class="text-muted-foreground text-sm">Registreret</p>
		{:else}
			<p class="text-lg font-semibold">
				{completedPercentage}%
			</p>
			<p class="text-muted-foreground text-sm">Udførsel</p>
		{/if}
	</div>
	<div class="flex flex-col items-center justify-center">
		{#if usePercent}
			<p class="text-lg font-semibold">
				{averagePercentage}%
			</p>
			<p class="text-muted-foreground text-sm">Gns / Registrering</p>
		{:else}
			<p class="text-lg font-semibold">
				{toDurationString(averageTime)}
			</p>
			<p class="text-muted-foreground text-sm">Gns / Registrering</p>
		{/if}
	</div>
</div>
