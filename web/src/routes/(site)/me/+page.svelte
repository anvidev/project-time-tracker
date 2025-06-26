<script lang="ts">
	import * as Navbar from '$lib/components/navbar';
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import { ArrowLeft } from '@lucide/svelte';
	import type { PageProps } from './$types';
	import CatTree from './CatTree.svelte';
	import CreateCategoryModal from './CreateCategoryModal.svelte';

	const { data }: PageProps = $props();

	const categoryTrees = $derived(data.categoryTrees);
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

	<div class="grid grid-cols-3 w-full gap-6 p-6">
		{#each categoryTrees as tree (tree.id)}
			<CatTree {tree} defaultOpen />
		{/each}
		<CreateCategoryModal
			triggerClass={buttonVariants({
				variant: 'outline',
				class: 'h-[42px] w-full cursor-pointer border-dashed'
			})}
			parentId={null}
			parentName={null}
		/>
	</div>
</section>
