<script lang="ts">
	import { superForm } from 'sveltekit-superforms';
	import type { PageProps } from './$types';
	import { Label } from '$lib/components/ui/label';
	import { Button } from '$lib/components/ui/button';
	import { DateFormatter } from '@internationalized/date';
	import Input from '$lib/components/ui/input/input.svelte';
	import type { Category, TimeEntry } from '$lib/types';
	import {
		Select,
		SelectContent,
		SelectGroup,
		SelectItem,
		SelectLabel,
		SelectTrigger
	} from '$lib/components/ui/select';
	import { Hour } from '$lib/utils';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';

	let { data }: PageProps = $props();

	const { form, constraints, enhance } = superForm(data.form, { dataType: 'json' });

	let categoryMap = $derived.by(() => {
		const res: Record<string, Category[]> = {};
		for (const category of data.categories) {
			if (res[category.rootTitle] == undefined) {
				res[category.rootTitle] = [];
			}

			res[category.rootTitle].push(category);
		}

		return res;
	});

	const daySummary = $derived.by(() => {
		data.daySummary.timeEntries.forEach((entry) => console.log(entry.duration));
		const grouped: Record<number, TimeEntry> = {};
		for (const entry of data.daySummary.timeEntries) {
			if (grouped[entry.categoryId] == undefined) {
				grouped[entry.categoryId] = {
					...entry
				};
			} else {
				grouped[entry.categoryId].duration += entry.duration;
			}
		}

		return {
			...data.daySummary,
			timeEntries: Array.from(Object.values(grouped)).sort((a, b) => b.duration - a.duration)
		};
	});

	const categoryTriggerContent = $derived(
		data.categories.find((c) => c.id == $form.categoryId)?.title ?? 'Vælg en kategori'
	);

	const remainingHours = $derived(
		Math.max((daySummary.maxHours - daySummary.totalHours) / Hour, 0)
	);

	const progressPercentage = $derived(
		Math.min((daySummary.totalHours / daySummary.maxHours) * 100, 100)
	);

	const getProgressColor = () => {
		if (progressPercentage < 50) return 'stroke-blue-500';
		if (progressPercentage < 75) return 'stroke-yellow-500';
		if (progressPercentage < 100) return 'stroke-orange-500';
		return 'stroke-green-500';
	};
</script>

<div class="grid w-full grid-cols-2 gap-2 p-2">
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
						{(daySummary.totalHours / Hour).toFixed(1)}t
					</div>
					<div class="text-sm text-slate-500">
						af {(daySummary.maxHours / Hour).toFixed(1)} timer
					</div>
					{#if remainingHours > 0}
						<div class="mt-1 text-xs text-slate-400">Mangler {remainingHours.toFixed(1)}t</div>
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
						{((daySummary.totalHours / daySummary.maxHours) * 100).toFixed(1)}%
					</p>
					<p class="text-muted-foreground text-sm">Udførsel</p>
				</div>
				<div class="flex flex-col items-center justify-center">
					<p class="text-lg font-semibold">
						{daySummary.totalHours / Hour / daySummary.timeEntries.length}t
					</p>
					<p class="text-muted-foreground text-sm">Gns / Registrering</p>
				</div>
			</div>
		</CardContent>
	</Card>

	<Card>
		<CardContent class="h-full">
			<form method="POST" class="flex h-full flex-col justify-evenly gap-4" use:enhance>
				<div class="grid gap-1">
					<Label class="gap-[2px]" for="category">Kategori<span class="text-red-700">*</span></Label
					>
					<Select
						type="single"
						name="category"
						value={$form.categoryId.toString()}
						onValueChange={(id) => ($form.categoryId = parseInt(id))}
					>
						<SelectTrigger class="w-full">
							{categoryTriggerContent}
						</SelectTrigger>
						<SelectContent>
							{#each Object.entries(categoryMap) as [label, categories]}
								<SelectGroup>
									<SelectLabel>{label}</SelectLabel>
									{#each categories as category}
										<SelectItem value={category.id.toString()}>
											{category.title}
										</SelectItem>
									{/each}
								</SelectGroup>
							{/each}
						</SelectContent>
					</Select>
				</div>

				<div class="grid gap-1">
					<Label class="gap-[2px]" for="duration">Timer<span class="text-red-700">*</span></Label>
					<Input
						name="duration"
						bind:value={$form.durationHours}
						type="number"
						step="0.01"
						placeholder="Indtast antal timer"
						{...$constraints.durationHours}
					/>
				</div>

				<div class="grid gap-1">
					<Label for="description">Beskrivelse</Label>
					<Input
						name="description"
						bind:value={$form.description}
						placeholder="Indtast valgfri beskrivelse"
						{...$constraints.description}
					/>
				</div>

				<Button type="submit" class="mt-2">Opret</Button>
			</form>
		</CardContent>
	</Card>

	<Card class="col-span-2">
		<CardHeader>
			<CardTitle>
				Registreringer d. {new DateFormatter('da-DK', { dateStyle: 'long' }).format(
					new Date(daySummary.date)
				)}
			</CardTitle>
			<CardDescription
				>Timer: {daySummary.totalHours / Hour} Max timer: {daySummary.maxHours /
					Hour}</CardDescription
			>
		</CardHeader>
		<CardContent>
			{#each daySummary.timeEntries as entry (entry.id)}
				<div class="bg-muted mb-2 flex items-center justify-between rounded-lg border p-2">
					<div>
						<p class="font-semibold tracking-tight">{entry.category}</p>
						<p class="text-muted-foreground text-sm">
							{((entry.duration / daySummary.totalHours) * 100).toFixed(1)}% af total
						</p>
					</div>
					<div>
						<Badge class="bg-background" variant="outline">
							{(entry.duration / Hour).toFixed(2)}t
						</Badge>
					</div>
				</div>
			{/each}
		</CardContent>
	</Card>
</div>
