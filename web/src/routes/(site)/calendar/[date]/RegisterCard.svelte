<script lang="ts">
	import { Card, CardContent } from '$lib/components/ui/card';
	import { cn } from '$lib/utils';
	import type { Category, LiveTimer } from '$lib/types';
	import TimeEntryForm from './TimeEntryForm.svelte';
	import LiveTimerForm from './LiveTimerForm.svelte';
	import type { ComponentProps } from 'svelte';

	const {
		timer,
		categories,
		showLive,
		...formProps
	}: {
		timer: LiveTimer | null;
		categories: Category[];
		// the live timer tracks time right now, so it is only offered for today (or while a timer is active)
		showLive: boolean;
	} & Omit<ComponentProps<typeof TimeEntryForm>, 'categories'> = $props();

	// both views stay mounted, so what is typed in one of them is not lost when switching
	let mode = $state<'manual' | 'live'>(timer ? 'live' : 'manual');
	const activeMode = $derived(showLive ? mode : 'manual');
</script>

{#snippet tab(value: 'manual' | 'live', text: string)}
	<button
		type="button"
		class={cn(
			'text-muted-foreground z-10 flex cursor-pointer items-center justify-center gap-2 text-sm outline-none',
			activeMode == value && 'text-primary-foreground font-semibold'
		)}
		onclick={() => (mode = value)}
	>
		{text}
		{#if value == 'live' && timer}
			<span
				class={cn(
					'size-2 rounded-full',
					timer.status == 'running' ? 'animate-pulse bg-green-500' : 'bg-amber-500'
				)}
			></span>
		{/if}
	</button>
{/snippet}

<Card>
	<CardContent class="flex h-full flex-col gap-4">
		{#if showLive}
			<div class="relative grid h-[36px] w-full grid-cols-2 items-center rounded-full border shadow-sm">
				<span
					class={cn(
						'bg-primary absolute z-0 h-[36px] w-1/2 rounded-full border shadow-sm transition-all',
						activeMode == 'live' && 'translate-x-full'
					)}
				></span>
				{@render tab('manual', 'Manuelt')}
				{@render tab('live', 'Live')}
			</div>
		{/if}

		<div class={cn('flex-1', activeMode != 'manual' && 'hidden')}>
			<TimeEntryForm {categories} {...formProps} />
		</div>
		{#if showLive}
			<div class={cn('flex-1', activeMode != 'live' && 'hidden')}>
				<LiveTimerForm {timer} {categories} />
			</div>
		{/if}
	</CardContent>
</Card>
