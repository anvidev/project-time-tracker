<script lang="ts">
	import { buttonVariants } from '$lib/components/ui/button';
	import { ArrowLeft, ArrowRight } from '@lucide/svelte';
	import { addMonths, format, subMonths } from 'date-fns';

	const { date }: { date: Date } = $props();

	const prevMonth = $derived(format(subMonths(date, 1), 'yyyy-MM-dd'));
	const nextMonth = $derived(format(addMonths(date, 1), 'yyyy-MM-dd'));
</script>

<div class="flex gap-2">
	<a class={`${buttonVariants({ variant: 'outline', size: 'icon' })}`} href={`?date=${prevMonth}`}>
		<ArrowLeft class="size-4" />
	</a>
	<a class={`${buttonVariants({ variant: 'outline', size: 'icon' })}`} href={`?date=${nextMonth}`}>
		<ArrowRight class="size-4" />
	</a>
</div>
