<script lang="ts">
	import * as Navbar from '$lib/components/navbar';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { ArrowLeft, Loader2 } from '@lucide/svelte';
	import type { PageProps } from './$types';
	import CatTree from './CatTree.svelte';
	import CreateCategoryModal from './CreateCategoryModal.svelte';
	import MaxHoursOverview from './MaxHoursOverview.svelte';
	import { superForm } from 'sveltekit-superforms';

	const { data }: PageProps = $props();

	const categoryTrees = $derived(data.categoryTrees);

	const maxHoursSuperForm = superForm(
		{ maxHours: data.maxHours },
		{
			dataType: 'json',
			delayMs: 350,
			timeoutMs: 8000,
			resetForm: false
		}
	);

	const { isTainted, tainted, submit, submitting, delayed } = maxHoursSuperForm;
</script>

<Navbar.Root>
	<Navbar.Action>
		<Button href="/calendar" variant="link">
			<ArrowLeft />
			Tilbage
		</Button>
	</Navbar.Action>

	<Navbar.Title>Profil</Navbar.Title>
</Navbar.Root>

<section class="flex w-full flex-col rounded-lg border">
	<div class="grid size-full grid-rows-2 border-b p-6 pb-4">
		<p class="row-start-1 w-full font-semibold tracking-tight">Kategorier</p>
		<p class="text-muted-foreground w-full text-xs">Se og opdater dine valgte kategorier</p>
	</div>

	<div class="grid w-full grid-cols-3 gap-6 p-6">
		{#each categoryTrees as tree (tree.id)}
			<CatTree {tree} defaultOpen />
		{/each}
		<CreateCategoryModal
			triggerClass={buttonVariants({
				variant: 'ghost',
				class: 'h-[41px] w-full cursor-pointer border border-dashed'
			})}
			parentId={null}
			parentName={null}
		/>
	</div>
</section>

<section class="flex w-full flex-col rounded-lg border">
	<div class="grid size-full grid-rows-2 border-b p-6 pb-4">
		<p class="row-start-1 w-full font-semibold tracking-tight">Timer</p>
		<p class="text-muted-foreground w-full text-xs">Se og opdater dine daglige timer</p>
		<Button
			class="row-span-2 row-start-1 my-auto ml-auto cursor-pointer"
			variant="outline"
			form="max-hours-form"
			onclick={() => submit()}
			disabled={!isTainted($tainted) || $submitting}
		>
			{#if $delayed}
				<Loader2 class="animate-spin" />
			{/if}
			Opdater
		</Button>
	</div>

	<MaxHoursOverview superForm={maxHoursSuperForm} />
</section>
