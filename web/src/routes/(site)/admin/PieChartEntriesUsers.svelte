<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { PieChart, Text } from 'layerchart';
	import type { TimeEntry } from '$lib/types';
	import { chartColors } from '$lib/chart-colors';
	import { toDurationString } from '$lib/utils';

	const { entries }: { entries: TimeEntry[] } = $props();

	const chartData = $derived(
		Array.from(
			entries
				.reduce((acc, cur) => {
					const user = cur.userName;

					if (!acc.has(user)) {
						acc.set(user, 0);
					}

					const oldOccurrence = acc.get(user);
					acc.set(user, (oldOccurrence ?? 0) + cur.duration);

					return acc;
				}, new Map<string, number>())
				.entries()
				.map(([key, value], index) => ({
					user: key,
					count: value,
					color: chartColors[index % chartColors.length]
				}))
		)
	);

	const chartConfig = $derived(
		entries.reduce(
			(acc, cur) => {
				const user = cur.userName;

				if (!acc[user]) {
					acc[user] = {};
				}

				acc[user] = {
					label: user
				};

				return acc;
			},
			{} as Record<string, any>
		)
	) satisfies Chart.ChartConfig;
</script>

<Card.Root class="flex flex-col">
	<Card.Header class="items-center">
		<Card.Title class="text-center">Timer fordelt på brugere</Card.Title>
	</Card.Header>
	<Card.Content class="flex-1">
		<Chart.Container config={chartConfig} class="mx-auto aspect-square max-h-[250px]">
			<PieChart
				data={chartData}
				innerRadius={65}
				key="user"
				padAngle={0.01}
				cornerRadius={6}
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
						value={String(chartData.length)}
						textAnchor="middle"
						verticalAnchor="middle"
						class="fill-foreground text-xl! font-bold tabular-nums"
						dy={-2}
					/>
					<Text
						value="Brugere"
						textAnchor="middle"
						verticalAnchor="middle"
						class="fill-muted-foreground! text-muted-foreground"
						dy={16}
					/>
				{/snippet}
				{#snippet tooltip()}
					<Chart.Tooltip
						valueClassName="font-semibold"
						hideLabel
						valueFormatter={(v) => toDurationString(v)}
					/>
				{/snippet}
			</PieChart>
		</Chart.Container>
	</Card.Content>
</Card.Root>
