<script lang="ts">
	import { cn } from '$lib/utils';

	let {
		value: activeItem = $bindable(),
		onActiveChange
	}: {
		value?: number | boolean;
		onActiveChange?: (active: 'left' | 'right') => void;
	} = $props();

	const toggle = () => {
		activeItem = Number(!Boolean(activeItem));
		if (onActiveChange) {
			onActiveChange(activeItem ? 'right' : 'left');
		}
	};
</script>

{#snippet element(text: string, isActive: boolean)}
	<span
		class={cn(
			'text-muted-foreground z-10 text-center text-sm',
			isActive && 'text-primary-foreground font-semibold'
		)}
	>
		{text}
	</span>
{/snippet}

<button
	class="relative grid h-[36px] w-[72px] cursor-pointer grid-cols-2 items-center rounded-xl border shadow-sm outline-none"
	onclick={toggle}
>
	<span
		class={`bg-primary absolute z-0 size-[36px] rounded-full border shadow-sm transition-all ${activeItem ? 'translate-x-[35px]' : ''}`}
	></span>
	{@render element('t', !Boolean(activeItem))}
	{@render element('%', Boolean(activeItem))}
</button>
