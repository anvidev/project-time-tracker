<script lang="ts">
	import ProgressContent from '$lib/components/progress-content.svelte';
	import { Hour, maxFractionDigits, toDurationString } from '$lib/utils';
	import { onMount } from 'svelte';
	import type { LayoutProps } from './$types';
	import { Separator } from '$lib/components/ui/separator';

	const entries = [2 * Hour, 1.5 * Hour, 1 * Hour, 1.5 * Hour, 2 * Hour];

	const { children }: LayoutProps = $props();
	const usePercent = false;
	let totalHours = $state(0 * Hour);
	let entryCount = $state(0);
	const maxHours = 8 * Hour;

	const totalHoursStr = $derived(toDurationString(totalHours));
	const maxHoursStr = $derived(toDurationString(maxHours));
	const remainingHours = $derived(maxHours - totalHours);
	const averageTime = $derived(
		entryCount && totalHours ? maxFractionDigits(totalHours / entryCount, 2) : 0
	);

	const completedPercentage = $derived(maxFractionDigits((totalHours / maxHours) * 100, 2));
	const averagePercentage = $derived(maxFractionDigits(completedPercentage / entryCount, 2));
	const progressPercentage = $derived(
		maxFractionDigits(Math.min((totalHours / maxHours) * 100, 100), 2)
	);

	const timeoutCallback = (count = 0) => {
		if (count < entries.length) {
			totalHours += entries[count++];
			entryCount = count;
			setTimeout(() => timeoutCallback(count), 2500);
		}
	};

	onMount(() => {
		setTimeout(timeoutCallback, 1000);
	});
</script>

<main class="relative grid min-h-svh lg:grid-cols-2">
	<div class="bg-muted hidden lg:grid lg:items-center lg:justify-center">
		<div
			class="bg-background grid items-center justify-center justify-items-center gap-12 rounded-xl px-12 py-8"
		>
			<div class="grid w-[40ch] gap-4 pt-4 text-center">
				<h1
					class="text-muted-foreground text-center text-2xl leading-none font-bold tracking-tight"
				>
					Velkommen til SkanCode Tid
				</h1>
				<p class="text-foreground/45 text-center text-sm leading-none">
					Simpel, nem og hurtig tidsregistrering
				</p>
			</div>
			<Separator />
			<div class="flex flex-col items-center gap-8">
				<ProgressContent
					{usePercent}
					{completedPercentage}
					{progressPercentage}
					{averagePercentage}
					{averageTime}
					{remainingHours}
					{totalHoursStr}
					{maxHoursStr}
					{entryCount}
				/>
			</div>
		</div>
	</div>
	<div class="m-auto flex flex-col p-6 md:p-10">
		{@render children()}
	</div>
</main>
