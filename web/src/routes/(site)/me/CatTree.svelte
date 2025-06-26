<script lang="ts">
	import Self from './CatTree.svelte';
	import type { CategoryTree } from '$lib/types';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import { ChevronsUpDown, Circle, CircleCheck, CircleDot } from '@lucide/svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { onMount } from 'svelte';

	const {
		tree,
		level = 0,
		parentIsFollowed = false
	}: {
		tree: CategoryTree;
		level?: number;
		parentIsFollowed?: boolean;
	} = $props();

	const isBrowser = typeof window !== 'undefined';

	let pixelsToLastChild = $state(0);

	const resizeObserver = $derived.by(() => {
		if (!isBrowser) {
			return null;
		}

		return new ResizeObserver((entries) => {
			const parentBounds = entries[0].target.firstElementChild?.getBoundingClientRect();

			const children = Array.from(
				document.querySelectorAll(`[data-parent-id='${tree.id}']`).values()
			);
			const lastChild = children.at(-1);
			const childBounds = lastChild?.getBoundingClientRect();

			if (parentBounds?.bottom && childBounds?.top) {
				pixelsToLastChild = childBounds.top - parentBounds.bottom;
			} else {
				pixelsToLastChild = 0;
			}
		});
	});

	const isNestedChildFollowed = (tree: CategoryTree): boolean => {
		return tree.children.some((child) => child.isFollowed || isNestedChildFollowed(child));
	};

	onMount(() => {
		const elem = document.querySelector(`[data-tree-id='${tree.id}']`);

		if (elem && resizeObserver) resizeObserver.observe(elem);

		return () => {
			resizeObserver?.disconnect();
		};
	});
</script>

<Collapsible.Root
	class="relative flex flex-col gap-1 transition-all"
	data-tree-id={`${tree.id}`}
	data-parent-id={`${tree.parentId}`}
>
	<div
		class="bg-background z-10 flex w-[350px] items-center gap-1 rounded-lg border px-2 py-1 text-sm tabular-nums"
	>
		{#if parentIsFollowed || tree.isFollowed}
			<CircleCheck class="size-4 text-green-600" />
		{:else if isNestedChildFollowed(tree)}
			<CircleDot class="size-4 text-yellow-600" />
		{:else}
			<Circle class="size-4" />
		{/if}
		<p>{tree.title}</p>
		<Collapsible.Trigger
			disabled={tree.children.length == 0}
			class={buttonVariants({ variant: 'ghost', size: 'sm', class: 'ml-auto' })}
		>
			<ChevronsUpDown />
		</Collapsible.Trigger>
	</div>
	{#if tree.children.length > 0}
		<Collapsible.Content>
			<div class="relative flex flex-col gap-1 pl-[32px]">
				<div
					class={`bg-border absolute w-px -translate-x-[16px] -translate-y-1`}
					style={`height: ${pixelsToLastChild + 4}px`}
				></div>
				{#each tree.children as child}
					<Self
						tree={child}
						level={level + 1}
						parentIsFollowed={parentIsFollowed || tree.isFollowed}
					/>
				{/each}
			</div>
		</Collapsible.Content>
	{/if}
	{#if level > 0}
		<div
			class="absolute top-[21px] z-0 h-[16px] w-[16px] -translate-x-full -translate-y-full rounded-bl-sm border-b border-l"
		></div>
	{/if}
</Collapsible.Root>
