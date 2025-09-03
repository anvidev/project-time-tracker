<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { toDurationString } from '$lib/utils';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { PieChart, Text } from 'layerchart';
	import type { TimeEntry } from '$lib/types';
	import { chartColors } from '$lib/chart-colors';

	const { entries, timeSpent }: { entries: TimeEntry[]; timeSpent: number } = $props();

	const chartData = $derived(
		Array.from(
			entries
				.reduce((acc, cur) => {
					const category = cur.category;
					const duration = cur.duration;

					if (!acc.has(category)) {
						acc.set(category, 0);
					}

					const oldDuration = acc.get(category);
					acc.set(category, (oldDuration ?? 0) + duration);

					return acc;
				}, new Map<string, number>())
				.entries()
				.map(([key, value], index) => ({
					category: key,
					hours: value,
					color: chartColors[index % chartColors.length]
				}))
		)
	);

	const chartConfig = $derived(
		entries.reduce(
			(acc, cur) => {
				const category = cur.category;

				if (!acc[category]) {
					acc[category] = {};
				}

				acc[category] = {
					label: category
				};

				return acc;
			},
			{} as Record<string, any>
		)
	) satisfies Chart.ChartConfig;

	let totalHours = $derived(timeSpent);
</script>

<Card.Root class="flex flex-col">
	<Card.Header class="items-center">
		<Card.Title class="text-center">Timer fordelt på kategorier</Card.Title>
	</Card.Header>
	<Card.Content class="flex-1">
		<Chart.Container config={chartConfig} class="mx-auto aspect-square max-h-[250px]">
			<PieChart
				data={chartData}
				innerRadius={65}
				padAngle={0.01}
				cornerRadius={6}
				key="category"
				value="hours"
				cRange={chartData.map((d) => d.color)}
				c="color"
				props={{
					pie: {
						motion: 'tween'
					}
				}}
			>
				{#snippet aboveMarks()}
					<Text
						value={toDurationString(totalHours)}
						textAnchor="middle"
						verticalAnchor="middle"
						class="fill-foreground text-xl! font-bold tabular-nums"
						dy={-2}
					/>
					<Text
						value="Brugt i alt"
						textAnchor="middle"
						verticalAnchor="middle"
						class="fill-muted-foreground! text-muted-foreground"
						dy={16}
					/>
				{/snippet}
				{#snippet tooltip()}
					<Chart.Tooltip
						valueFormatter={(v) => toDurationString(v)}
						valueClassName="font-semibold"
						hideLabel
					/>
				{/snippet}
			</PieChart>
		</Chart.Container>
	</Card.Content>
</Card.Root>
