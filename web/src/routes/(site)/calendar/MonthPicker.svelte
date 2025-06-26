<script lang="ts">
	import { Button, buttonVariants } from '$lib/components/ui/button';
	import * as Select from '$lib/components/ui/select';
	import { DateFormatter, getLocalTimeZone, today } from '@internationalized/date';
	import { ArrowLeft, ArrowRight } from '@lucide/svelte';
	import { addMonths, format, startOfMonth, subMonths } from 'date-fns';
	import { months } from '$lib/types';
	import { goto } from '$app/navigation';

	const { date }: { date: Date } = $props();

	const prevMonth = $derived(subMonths(startOfMonth(date), 1));
	const prevMonthStr = $derived(format(prevMonth, 'yyyy-MM-dd'));

	const nextMonth = $derived(addMonths(startOfMonth(date), 1));
	const nextMonthStr = $derived(format(nextMonth, 'yyyy-MM-dd'));

	let monthValue = $derived(date.getMonth());
	let yearValue = $derived(date.getFullYear());

	const onMonthChange = (month: string) => {
		monthValue = parseInt(month);
		const newDate = new Date(yearValue, monthValue, 1);
		goto(`?date=${format(newDate, 'yyyy-MM-dd')}`);
	};
	const onYearChange = (year: string) => {
		yearValue = parseInt(year);
		const newDate = new Date(yearValue, monthValue, 1);
		goto(`?date=${format(newDate, 'yyyy-MM-dd')}`);
	};

	const maxYear = today(getLocalTimeZone()).year + 1;
	const years = Array.from({ length: maxYear - 2025 + 1 }, (_, i) => 2025 + i);

	const monthDF = new DateFormatter('da-DK', {
		month: 'long'
	});
	const yearDF = new DateFormatter('da-DK', {
		year: 'numeric'
	});
</script>

<div class="flex gap-2">
	{#if prevMonth.getFullYear() < 2025}
		<Button variant="outline" size="icon" disabled>
			<ArrowLeft class="size-4" />
		</Button>
	{:else}
		<a
			class={`${buttonVariants({ variant: 'outline', size: 'icon' })}`}
			href={`?date=${prevMonthStr}`}
		>
			<ArrowLeft class="size-4" />
		</a>
	{/if}
	<div class="md:hidden flex items-center justify-center text-sm font-medium capitalize">
		{monthDF.format(new Date(yearValue, monthValue, 1))} {yearDF.format(new Date(yearValue, monthValue, 1))}
	</div>
	<Select.Root type="single" value={monthValue.toString()} onValueChange={onMonthChange}>
		<Select.Trigger class="cursor-pointer capitalize hidden md:flex">
			{monthDF.format(new Date(yearValue, monthValue, 1))}
		</Select.Trigger>
		<Select.Content>
			{#each months as _, index}
				<Select.Item class="cursor-pointer capitalize" value={index.toString()}>
					{monthDF.format(new Date(yearValue, index, 1))}
				</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>
	<Select.Root type="single" value={yearValue.toString()} onValueChange={onYearChange}>
		<Select.Trigger class="cursor-pointer capitalize hidden md:flex">
			{yearDF.format(new Date(yearValue, monthValue, 1))}
		</Select.Trigger>
		<Select.Content>
			{#each years as year}
				<Select.Item class="cursor-pointer capitalize" value={year.toString()}>
					{yearDF.format(new Date(year, monthValue, 1))}
				</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>
	{#if nextMonth.getFullYear() > maxYear}
		<Button variant="outline" size="icon" disabled>
			<ArrowRight class="size-4" />
		</Button>
	{:else}
		<a
			class={`${buttonVariants({ variant: 'outline', size: 'icon' })}`}
			href={`?date=${nextMonthStr}`}
		>
			<ArrowRight class="size-4" />
		</a>
	{/if}
</div>
