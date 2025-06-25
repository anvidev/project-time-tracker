<script lang="ts">
	import type { CategoryTree } from '$lib/types';
	import type { PageProps } from './$types';
	import CatTree from './CatTree.svelte';

	const { data }: PageProps = $props();
</script>

{#snippet CategoryTree(tree: CategoryTree, level: number = 0, parentIsFollowed: boolean = false)}
	<CatTree {tree} {level} {parentIsFollowed} />
{/snippet}

<div class="grid w-full gap-2">
	{#each data.categoryTrees as tree (tree.id)}
		{@render CategoryTree(tree)}
	{/each}
</div>
