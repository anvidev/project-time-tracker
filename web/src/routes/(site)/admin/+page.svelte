<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { cn, toDurationString } from '$lib/utils';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Eye } from '@lucide/svelte';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { PieChart, Text, BarChart, Highlight, type ChartContextValue } from 'layerchart';
	import { scaleBand } from 'd3-scale';
	import { cubicInOut } from 'svelte/easing';

	const { data } = $props();
	const { entries } = $derived(data);

	const pieCountData = $derived(
		Array.from(
			entries.entries
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
				.map(([key, value]) => ({ category: key, count: value, color: 'var(--chart-3)' }))
		)
	);

	const pieCountConfig = $derived(
		entries.entries.reduce(
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

	let totalHours = $derived(entries.timeSpent);
	const pieHoursData = $derived(
		Array.from(
			entries.entries
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
				.map(([key, value]) => ({ category: key, hours: value, color: 'var(--chart-3)' }))
		)
	);

	const pieHoursConfig = $derived(
		entries.entries.reduce(
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

	const pieUsersData = $derived(
		Array.from(
			entries.entries
				.reduce((acc, cur) => {
					const user = cur.userName;

					if (!acc.has(user)) {
						acc.set(user, 0);
					}

					const oldOccurrence = acc.get(user);
					acc.set(user, (oldOccurrence ?? 0) + 1);

					return acc;
				}, new Map<string, number>())
				.entries()
				.map(([key, value]) => ({ user: key, count: value, color: 'var(--chart-3)' }))
		)
	);

	const pieUsersConfig = $derived(
		entries.entries.reduce(
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

	const chartData = [
		{ month: 'January', desktop: 186, mobile: 80 },
		{ month: 'February', desktop: 305, mobile: 200 },
		{ month: 'March', desktop: 237, mobile: 120 },
		{ month: 'April', desktop: 73, mobile: 190 },
		{ month: 'May', desktop: 209, mobile: 130 },
		{ month: 'June', desktop: 214, mobile: 140 },
		{ month: 'July', desktop: 214, mobile: 140 },
		{ month: 'August', desktop: 214, mobile: 140 },
		{ month: 'September', desktop: 214, mobile: 140 },
		{ month: 'October', desktop: 214, mobile: 140 },
		{ month: 'November', desktop: 214, mobile: 140 },
		{ month: 'December', desktop: 214, mobile: 140 }
	];

	const chartConfig = {
		desktop: { label: 'Desktop', color: 'var(--chart-3)' },
		mobile: { label: 'Mobile', color: 'var(--chart-3)' }
	} satisfies Chart.ChartConfig;

	let context = $state<ChartContextValue>();
</script>

<h1>admin</h1>

<div class="flex w-full flex-col gap-4 lg:flex-row lg:items-center">
	<div class="w-1/3">
		<Card.Root class="flex flex-col">
			<Card.Header class="items-center">
				<Card.Title>Antal tidsregistreringer på kategorier</Card.Title>
				<Card.Description>January - June 2024</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={pieCountConfig} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={pieCountData}
						innerRadius={65}
						padAngle={0.01}
						cornerRadius={6}
						key="category"
						value="count"
						cRange={pieCountData.map((d) => d.color)}
						c="color"
						props={{
							pie: {
								motion: 'tween'
							}
						}}
					>
						{#snippet aboveMarks()}
							<Text
								value={String(pieCountData.reduce((acc, cur) => acc + cur.count, 0))}
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
	</div>

	<div class="w-1/3">
		<Card.Root class="flex flex-col">
			<Card.Header class="items-center">
				<Card.Title>Timer brugt på kategorier</Card.Title>
				<Card.Description>January - June 2024</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={pieHoursConfig} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={pieHoursData}
						innerRadius={65}
						padAngle={0.01}
						cornerRadius={6}
						key="category"
						value="hours"
						cRange={pieCountData.map((d) => d.color)}
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
	</div>

	<div class="w-1/3">
		<Card.Root class="flex flex-col">
			<Card.Header class="items-center">
				<Card.Title>Tidsregistreringer fordelt på brugere</Card.Title>
				<Card.Description>January - June 2024</Card.Description>
			</Card.Header>
			<Card.Content class="flex-1">
				<Chart.Container config={pieUsersConfig} class="mx-auto aspect-square max-h-[250px]">
					<PieChart
						data={pieUsersData}
						innerRadius={65}
						key="user"
						padAngle={0.01}
						cornerRadius={6}
						value="count"
						cRange={pieUsersData.map((d) => d.color)}
						c="color"
						props={{
							pie: {
								motion: 'tween'
							}
						}}
					>
						{#snippet aboveMarks()}
							<Text
								value={String(pieUsersData.length)}
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
							<Chart.Tooltip valueClassName="font-semibold" hideLabel />
						{/snippet}
					</PieChart>
				</Chart.Container>
			</Card.Content>
		</Card.Root>
	</div>
</div>

<div class="w-full">
	<Card.Root>
		<Card.Header>
			<Card.Title>Bar Chart - Stacked + Legend</Card.Title>
			<Card.Description>January - June 2024</Card.Description>
		</Card.Header>
		<Card.Content>
			<Chart.Container config={chartConfig} class="max-h-[250px] w-full">
				<BarChart
					bind:context
					data={chartData}
					stackPadding={2}
					xScale={scaleBand().padding(0.25)}
					x="month"
					axis="x"
					rule={false}
					series={[
						{
							key: 'desktop',
							label: 'Desktop',
							color: chartConfig.desktop.color
						},
						{
							key: 'mobile',
							label: 'Mobile',
							color: chartConfig.mobile.color
						}
					]}
					seriesLayout="stack"
					props={{
						bars: {
							stroke: 'none',
							initialY: context?.height,
							initialHeight: 0,
							motion: {
								y: { type: 'tween', duration: 500, easing: cubicInOut },
								height: { type: 'tween', duration: 500, easing: cubicInOut }
							},
							radius: 6,
							rounded: 'all'
						},
						highlight: { area: false },
						xAxis: { format: (d) => d.slice(0, 3) }
					}}
				>
					{#snippet belowMarks()}
						<Highlight area={{ class: 'fill-muted' }} />
					{/snippet}

					{#snippet tooltip()}
						<Chart.Tooltip />
					{/snippet}
				</BarChart>
			</Chart.Container>
		</Card.Content>
	</Card.Root>
</div>

<div class="bg-background w-full rounded-md border shadow-sm">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head class="w-28">Dato</Table.Head>
				<Table.Head class="w-44">Bruger</Table.Head>
				<Table.Head class="w-44">Kategori</Table.Head>
				<Table.Head class="hidden max-w-80 lg:table-cell">Beskrivelse</Table.Head>
				<Table.Head class="w-44 text-right">Timer</Table.Head>
				<Table.Head class="w-20 text-right"></Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each entries.entries as entry (entry)}
				<Table.Row>
					<Table.Cell class="w-28">{entry.date}</Table.Cell>
					<Table.Cell class="w-44">{entry.userName}</Table.Cell>
					<Table.Cell class="w-44">{entry.category}</Table.Cell>
					<Table.Cell
						class={cn(
							'hidden max-w-80 truncate lg:table-cell',
							entry.description == '' && 'text-muted-foreground'
						)}>{entry.description != '' ? entry.description : 'Ingen beskrivelse'}</Table.Cell
					>
					<Table.Cell class="w-44 text-right">{toDurationString(entry.duration)}</Table.Cell>
					<Table.Cell class="w-20">
						<Eye class="ml-auto size-4" />
					</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
