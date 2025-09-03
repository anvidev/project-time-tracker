<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import type { TimeEntry } from '$lib/types';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { BarChart, Highlight, type ChartContextValue } from 'layerchart';
	import { scaleBand } from 'd3-scale';
	import { cubicInOut } from 'svelte/easing';
	import { chartColors } from '$lib/chart-colors';
	import { toDurationString } from '$lib/utils';

	const { entries }: { entries: TimeEntry[] } = $props();

	function transformData(registrations: TimeEntry[]) {
		const grouped = registrations.reduce(
			(acc, reg) => {
				const date = new Date(reg.date);
				const monthKey = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;
				const monthLabel = date.toLocaleDateString('da-DK', {
					year: 'numeric',
					month: 'long'
				});

				if (!acc[monthKey]) {
					acc[monthKey] = { month: monthLabel, monthKey, categories: {} };
				}

				if (!acc[monthKey].categories[reg.category]) {
					acc[monthKey].categories[reg.category] = 0;
				}

				acc[monthKey].categories[reg.category] =
					acc[monthKey].categories[reg.category] + reg.duration;

				return acc;
			},
			{} as Record<string, { month: string; monthKey: string; categories: Record<string, number> }>
		);

		const allCategories = [...new Set(registrations.map((reg) => reg.category))];

		const chartData = Object.values(grouped)
			.sort((a, b) => a.monthKey.localeCompare(b.monthKey))
			.map((item) => {
				const result: Record<string, any> = { month: item.month };
				allCategories.forEach((category) => {
					result[category] = item.categories[category] || 0;
				});
				return result;
			});

		return { chartData, categories: allCategories };
	}

	const { chartData, categories } = transformData(entries);

	const chartConfig = categories.reduce(
		(config, category, index) => {
			config[category] = {
				label: category,
				color: chartColors[index % chartColors.length]
			};
			return config;
		},
		{} as Record<string, { label: string; color: string }>
	);

	const series = categories.map((category, index) => ({
		key: category,
		label: category,
		color: chartConfig[category].color,
		props: index === 0 ? { rounded: 'bottom' } : {}
	}));

	let context = $state<ChartContextValue>();
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Registreringer på kategorier fordelt på måneder</Card.Title>
	</Card.Header>
	<Card.Content>
		<Chart.Container config={chartConfig} class="mx-auto max-h-[400px]">
			<BarChart
				bind:context
				data={chartData}
				xScale={scaleBand().padding(0.25)}
				x="month"
				axis="x"
				rule={false}
				stackPadding={2}
				{series}
				seriesLayout="stack"
				props={{
					bars: {
						radius: 4,
						rounded: 'all',
						stroke: 'none',
						initialY: context?.height,
						initialHeight: 0,
						motion: {
							y: { type: 'tween', duration: 500, easing: cubicInOut },
							height: { type: 'tween', duration: 500, easing: cubicInOut }
						}
					},
					highlight: { area: false },
					xAxis: { format: (d) => d.slice(0, 3) }
				}}
			>
				{#snippet belowMarks()}
					<Highlight area={{ class: 'fill-muted' }} />
				{/snippet}

				{#snippet tooltip()}
					<Chart.Tooltip valueFormatter={(v) => toDurationString(v)} />
				{/snippet}
			</BarChart>
		</Chart.Container>
	</Card.Content>
</Card.Root>
