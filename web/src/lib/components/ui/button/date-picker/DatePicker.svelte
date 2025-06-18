<script lang="ts">
	import CalendarIcon from '@lucide/svelte/icons/calendar';
	import { DateFormatter, type DateValue, getLocalTimeZone } from '@internationalized/date';
	import { cn } from '$lib/utils.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { Calendar } from '$lib/components/ui/calendar/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';

	let {
		value,
		placeholder = 'Pick a date',
		dateFormatter: df = new DateFormatter('da-DK', {
			dateStyle: 'long'
		}),
		onValueChange,
		class: clazz,
		name,
	}: {
		value?: DateValue;
		placeholder?: string;
		dateFormatter?: DateFormatter;
		onValueChange?: (v: DateValue | undefined) => void;
		class?: string;
		name?: string;
	} = $props();

	let contentRef = $state<HTMLElement | null>(null);
</script>

<Popover.Root>
	<Popover.Trigger
		class={cn(
			buttonVariants({
				variant: 'outline',
				class: cn('w-[280px] justify-start text-left font-normal', clazz)
			}),
			!value && 'text-muted-foreground'
		)}
		{name}
	>
		<CalendarIcon />
		{value ? df.format(value.toDate(getLocalTimeZone())) : placeholder}
	</Popover.Trigger>
	<Popover.Content bind:ref={contentRef} class="w-auto p-0">
		<Calendar type="single" bind:value {onValueChange} />
	</Popover.Content>
</Popover.Root>
