<script lang="ts">
	import type { SummaryDay } from '$lib/types';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Hour, maxFractionDigits } from '$lib/utils';

	const { daySummary }: { daySummary: SummaryDay } = $props();

	const averageTime = $derived.by(() => {
		if (daySummary.totalHours == 0 || daySummary.timeEntries.length == 0) {
			return 0;
		}

		return maxFractionDigits(daySummary.totalHours / Hour / daySummary.timeEntries.length, 2);
	});

	const totalHours = $derived(
		daySummary.totalHours > 0.0 ? maxFractionDigits(daySummary.totalHours / Hour, 1) : 0
	);
	const maxHours = $derived(
		daySummary.maxHours > 0.0 ? maxFractionDigits(daySummary.maxHours / Hour, 1) : 0
	);

	const completedPercentage = $derived(
		totalHours < 0.1 || maxHours < 0.1 ? 0 : maxFractionDigits((totalHours / maxHours) * 100, 2)
	);

	const remainingHours = $derived(
		totalHours < 0.1 || maxHours < 0.1
			? 0
			: maxFractionDigits(Math.max((maxHours - totalHours) / Hour, 0), 1)
	);

	const progressPercentage = $derived(
		totalHours < 0.1 || maxHours < 0.1
			? 0
			: maxFractionDigits(Math.min((totalHours / maxHours) * 100, 100), 2)
	);

	const getProgressColor = () => {
		if (progressPercentage < 25) return 'stroke-red-500';
		if (progressPercentage < 50) return 'stroke-orange-500';
		if (progressPercentage < 75) return 'stroke-yellow-500';
		if (progressPercentage < 100) return 'stroke-green-300';
		return 'stroke-green-500';
	};
</script>

<Card>
	<CardHeader class="text-center">
		<CardTitle class="text-2xl">Dagens opsummering</CardTitle>
	</CardHeader>
	<CardContent class="flex h-full flex-col items-center gap-8">
		<div class="relative h-48 w-48">
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
				<div class="text-3xl font-bold text-slate-800">
					{totalHours}t
				</div>
				<div class="text-sm text-slate-500">
					af {maxHours} timer
				</div>
				{#if remainingHours > 0}
					<div class="mt-1 text-xs text-slate-400">Mangler {remainingHours}t</div>
				{/if}
			</div>
		</div>
		<div class="grid w-full grid-cols-3">
			<div class="flex flex-col items-center justify-center">
				<p class="text-lg font-semibold">{daySummary.timeEntries.length}</p>
				<p class="text-muted-foreground text-sm">Registreringer</p>
			</div>
			<div class="flex flex-col items-center justify-center">
				<p class="text-lg font-semibold">
					{completedPercentage}%
				</p>
				<p class="text-muted-foreground text-sm">Udførsel</p>
			</div>
			<div class="flex flex-col items-center justify-center">
				<p class="text-lg font-semibold">
					{averageTime}t
				</p>
				<p class="text-muted-foreground text-sm">Gns / Registrering</p>
			</div>
		</div>
	</CardContent>
</Card>
