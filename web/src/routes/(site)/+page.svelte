<script lang="ts">
	import SuperDebug, { superForm } from 'sveltekit-superforms';
	import type { PageProps } from './$types';
	import { Label } from '$lib/components/ui/label';
	import DatePicker from '$lib/components/ui/button/date-picker/DatePicker.svelte';
	import { Button } from '$lib/components/ui/button';
	import {
		DateFormatter,
		fromDate,
		getLocalTimeZone,
		type DateValue
	} from '@internationalized/date';
	import { format, toDate } from 'date-fns';
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
	import { Separator } from '$lib/components/ui/separator';
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

	let dateValue = $derived(
		$form.date ? fromDate(toDate($form.date), getLocalTimeZone()) : undefined
	);

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
			timeEntries: Object.values(grouped)
		};
	});

	const categoryTriggerContent = $derived(
		data.categories.find((c) => c.id == $form.categoryId)?.title ?? 'Vælg en kategori'
	);

	const onValueChange = (v: DateValue | undefined) => {
		if (v) {
			$form.date = format(v.toDate(getLocalTimeZone()), 'yyyy-MM-dd');
		} else {
			$form.date = '';
		}
	};
</script>

<Card class="w-[350px]">
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
			<div class="mb-2 flex items-center justify-between rounded-lg border p-2">
				<div>
					<p class="font-semibold tracking-tight">{entry.category}</p>
					<p class="text-muted-foreground text-sm">
						{((entry.duration / daySummary.totalHours) * 100).toFixed(1)}% af total
					</p>
				</div>
				<div>
					<Badge variant="outline">{entry.duration / Hour}T</Badge>
				</div>
			</div>
		{/each}
	</CardContent>
</Card>

<Card class="w-[350px]">
	<CardContent>
		<form method="POST" class="flex flex-col gap-4" use:enhance>
			<div class="grid gap-1">
				<Label class="gap-[2px]">Dato<span class="text-red-700">*</span></Label>
				<DatePicker value={dateValue} {onValueChange} placeholder="Vælg en dato" class="w-full" />
			</div>

			<div class="grid gap-1">
				<Label class="gap-[2px]" for="category">Kategori<span class="text-red-700">*</span></Label>
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

<SuperDebug data={$form} />
