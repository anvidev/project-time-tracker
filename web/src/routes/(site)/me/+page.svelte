<script lang="ts">
	import * as Navbar from '$lib/components/navbar';
	import { Button } from '$lib/components/ui/button';
	import { ArrowLeft } from '@lucide/svelte';
	import type { PageProps } from './$types';
	import CatTree from './CatTree.svelte';

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
	<div class="size-full border-b p-6 pb-4">
		<p class="w-full text-lg font-semibold tracking-tight">Kategorier</p>
		<p class="text-muted-foreground w-full text-sm">Se og opdater dine valgte kategorier</p>
	</div>

	<div class="grid w-full gap-2 p-6">
		{#each categoryTrees as tree (tree.id)}
			<CatTree {tree} />
		{/each}
	</div>
</section>
