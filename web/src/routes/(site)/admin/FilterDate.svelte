<script lang="ts">
	import CalendarIcon from '@lucide/svelte/icons/calendar';
	import {
		type DateValue,
		CalendarDate,
		DateFormatter,
		getLocalTimeZone,
		parseDate
	} from '@internationalized/date';
	import { cn } from '$lib/utils.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Calendar } from '$lib/components/ui/calendar/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import Label from '$lib/components/ui/label/label.svelte';

	let {
		param = $bindable(),
		label,
		placeholder = 'Vælg en dato'
	}: {
		param: Date | null;
		label: string;
		placeholder?: string;
	} = $props();

	const df = new DateFormatter('da', {
		dateStyle: 'long'
	});

	let value = $state<DateValue | undefined>(
		param ? new CalendarDate(param.getFullYear(), param.getMonth() + 1, param.getDate()) : undefined
	);

	// Sync value when param changes externally
	$effect(() => {
		if (param) {
			value = new CalendarDate(param.getFullYear(), param.getMonth() + 1, param.getDate());
		} else {
			value = undefined;
		}
	});

	function onChangeFn(dateValue: DateValue | undefined) {
		value = dateValue;
		// Update the param prop using UTC to avoid timezone issues
		if (dateValue) {
			param = new Date(Date.UTC(dateValue.year, dateValue.month - 1, dateValue.day));
		} else {
			param = null;
		}
	}
</script>

<div class="grid gap-1.5">
	<Label>{label}</Label>
	<Popover.Root>
		<Popover.Trigger>
			{#snippet child({ props })}
				<Button
					variant="outline"
					class={cn(
						'w-[280px] justify-start text-left font-normal',
						!value && 'text-muted-foreground'
					)}
					{...props}
				>
					<CalendarIcon class="mr-2 size-4" />
					{value ? df.format(value.toDate(getLocalTimeZone())) : placeholder}
				</Button>
			{/snippet}
		</Popover.Trigger>
		<Popover.Content class="w-auto p-0">
			<Calendar locale="da" bind:value type="single" initialFocus onValueChange={onChangeFn} />
		</Popover.Content>
	</Popover.Root>
</div>
