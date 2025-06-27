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
  import { Textarea } from "$lib/components/ui/textarea/index.js";

	const {
		formData,
		categories,
		usePercent,
		maxHours
	}: {
		categories: Category[];
		usePercent: boolean;
		maxHours: number;
		formData: SuperValidated<
			{
				categoryId: number;
				date: string;
				durationHours: number | string;
				description?: string | undefined;
			},
			any,
			{
				categoryId: number;
				date: string;
				durationHours: number | string;
				description?: string | undefined;
			}
		>;
	} = $props();

	let percentState = $state(0);

	const { form, constraints, enhance, delayed } = superForm(formData, {
		dataType: 'json',
		delayMs: 150,
		timeoutMs: 8000,
		onUpdated: () => (percentState = 0)
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
					<Label class="gap-[2px]" for="duration">
						{#if usePercent}
							Procent
						{:else}
							Timer
						{/if}
						<span class="text-red-700">*</span></Label
					>
					{#if usePercent}
						<Input
							class="input-arrows-none"
							name="duration"
							bind:value={percentState}
							type="number"
							step="0.1"
							placeholder="Indtast procent"
							oninput={() => ($form.durationHours = maxHours * (percentState / 100))}
						/>
					{:else}
						<Input
							class="input-arrows-none"
							name="duration"
							bind:value={$form.durationHours}
							type="text"
							placeholder="Indtast antal timer"
							{...$constraints.durationHours}
							pattern="(?:\d+|\d+t \d+m|0t)"
						/>
					{/if}
				</div>

				<div class="grid gap-1">
					<Label for="description">Beskrivelse</Label>
					<Textarea
						name="description"
						bind:value={$form.description}
						placeholder="Indtast valgfri beskrivelse"
						class="resize-none"
						{...$constraints.description}
					></Textarea>
				</div>
			</div>

			<Button type="submit" class="mt-2 cursor-pointer">
				Opret
				{#if $delayed}
					<LoaderCircle class="animate-spin" />
				{/if}
			</Button>
		</form>
	</CardContent>
</Card>
