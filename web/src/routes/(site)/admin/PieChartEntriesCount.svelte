<script lang="ts">
	import type { TimeEntry } from '$lib/types';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { PieChart, Text } from 'layerchart';
	import { chartColors } from '$lib/chart-colors';

	const { entries }: { entries: TimeEntry[] } = $props();

	const chartData = $derived(
		Array.from(
			entries
				.reduce((acc, cur) => {
					const category = cur.category;

					if (!acc.has(category)) {
						acc.set(category, 0);
					}

					const occurence = acc.get(category);
					acc.set(category, (occurence ?? 0) + 1);

					return acc;
				}, new Map<string, number>())
				.entries()
				.map(([key, value], index) => ({
					category: key,
					count: value,
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
</script>

<Card.Root class="flex flex-col">
	<Card.Header class="items-center">
		<Card.Title class="text-center">Tidsregistreringer på kategorier</Card.Title>
	</Card.Header>
	<Card.Content class="flex-1">
		<Chart.Container config={chartConfig} class="mx-auto aspect-square max-h-[250px]">
			<PieChart
				data={chartData}
				innerRadius={65}
				padAngle={0.01}
				cornerRadius={6}
				key="category"
				value="count"
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
						value={String(chartData.reduce((acc, cur) => acc + cur.count, 0))}
						textAnchor="middle"
						verticalAnchor="middle"
						class="fill-foreground text-xl! font-bold tabular-nums"
						dy={-2}
					/>
					<Text
						value="Registreringer"
						textAnchor="middle"
						verticalAnchor="middle"
						class="fill-muted-foreground! text-muted-foreground"
						dy={16}
					/>
				{/snippet}
				{#snippet tooltip()}
					<Chart.Tooltip valueClassName="font-semibold" hideLabel />
				{/snippet}
			</PieChart>
		</Chart.Container>
	</Card.Content>
</Card.Root>
