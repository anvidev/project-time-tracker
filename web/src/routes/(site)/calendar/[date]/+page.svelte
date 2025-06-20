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
	import { ArrowLeft } from "@lucide/svelte";
	import ProgressCard from './ProgressCard.svelte';

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

	const dateFormatter = new DateFormatter('da-DK', { dateStyle: 'long' });

	const formattedDate = $derived(
		dateFormatter.format(
			new Date(daySummary.date)
		)
	)

	const categoryTriggerContent = $derived(
		data.categories.find((c) => c.id == $form.categoryId)?.title ?? 'Vælg en kategori'
	);
</script>

<div class="flex w-full items-center p-4 border rounded-xl mb-6 gap-4">
	<Button href="/calendar">
		<ArrowLeft />
		Tilbage
	</Button>
	<p class="font-semibold text-xl tracking-tight">Tidsregistrering d. {formattedDate}</p>
</div>

<div class="grid w-full grid-cols-2 gap-6">
	<ProgressCard {daySummary} />

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
				Registreringer d. {formattedDate}
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
