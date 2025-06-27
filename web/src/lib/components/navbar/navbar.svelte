<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { cn, type WithElementRef } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<div
	bind:this={ref}
	class={cn(
		'relative flex h-[70px] w-full items-center gap-4 rounded-lg border p-4 text-center shadow-xs',
		className
	)}
	{...restProps}
>
	{@render children?.()}
</div>
