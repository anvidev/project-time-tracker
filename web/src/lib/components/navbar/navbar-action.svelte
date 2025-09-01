<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { cn, type WithElementRef } from '$lib/utils.js';

	interface Props extends WithElementRef<HTMLAttributes<HTMLDivElement>> {
		side?: 'left' | 'right' 
	}

	let {
		ref = $bindable(null),
		class: className,
		children,
		side = 'left',
		...restProps
	}: Props = $props()


</script>

<div 
	bind:this={ref}
	class={cn(
		"absolute top-1/2 -translate-y-1/2 flex items-center gap-2",
		`${side}-4`,
		//side == 'left' ? 'left-4' : 'right-4',
		className,
	)}
	{...restProps}
>
	{@render children?.()}
</div>
