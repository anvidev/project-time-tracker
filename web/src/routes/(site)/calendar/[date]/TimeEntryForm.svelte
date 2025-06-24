<script lang="ts">
	import { superForm, type SuperValidated } from 'sveltekit-superforms';
	import {
		Select,
		SelectContent,
		SelectGroup,
		SelectItem,
		SelectLabel,
		SelectTrigger
	} from '$lib/components/ui/select';
	import { Card, CardContent } from '$lib/components/ui/card';
	import { Label } from '$lib/components/ui/label';
	import Input from '$lib/components/ui/input/input.svelte';
	import { Button } from '$lib/components/ui/button';
	import type { Category } from '$lib/types';
	import { LoaderCircle } from '@lucide/svelte';

	const {
		formData,
		categories
	}: {
		categories: Category[];
		formData: SuperValidated<
			{
				categoryId: number;
				date: string;
				durationHours: number;
				description?: string | undefined;
			},
			any,
			{
				categoryId: number;
				date: string;
				durationHours: number;
				description?: string | undefined;
			}
		>;
	} = $props();

	const { form, constraints, enhance, delayed } = superForm(formData, {
		dataType: 'json',
		delayMs: 150,
		timeoutMs: 8000,
	});

	let categoryMap = $derived.by(() => {
		const res: Record<string, Category[]> = {};
		for (const category of categories) {
			if (res[category.rootTitle] == undefined) {
				res[category.rootTitle] = [];
			}

			res[category.rootTitle].push(category);
		}

		return res;
	});
	const categoryTriggerContent = $derived(
		categories.find((c) => c.id == $form.categoryId)?.title ?? 'Vælg en kategori'
	);
</script>

<Card>
	<CardContent class="h-full">
		<form
			method="POST"
			action="?/createTimeEntry"
			class="flex h-full flex-col justify-between gap-4"
			use:enhance
		>
			<div class="flex flex-col gap-4">
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
			</div>

			<Button type="submit" class="mt-2 cursor-pointer">
				Opret
				{#if $delayed} <LoaderCircle class="animate-spin" /> {/if}
			</Button>
		</form>
	</CardContent>
</Card>
