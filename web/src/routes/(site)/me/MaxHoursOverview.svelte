<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import type { WeekdayHours } from '$lib/types';
	import { Minute, toDurationString } from '$lib/utils';
	import { ChevronDown, ChevronUp } from '@lucide/svelte';
	import { type SuperForm } from 'sveltekit-superforms';

	const {
		superForm: superFormProp
	}: {
		superForm: SuperForm<{ maxHours: WeekdayHours[] }>;
	} = $props();

	const { form, enhance, submitting } = superFormProp;

	const buttonUpClicked = (index: number) => {
		form.update(($form) => {
			$form.maxHours[index].hours += 5 * Minute;
			return $form;
		});
	};
	const buttonDownClicked = (index: number) => {
		form.update(
			($form) => {
				$form.maxHours[index].hours -= 5 * Minute;
				return $form;
			},
			{ taint: true }
		);
	};
</script>

<form id="max-hours-form" action="?/updateMaxHours" method="POST" use:enhance>
	<div class="grid w-full gap-6 p-6 lg:grid-cols-7">
		{#each $form.maxHours as hour, index (`${hour.weekday}_${toDurationString(hour.hours)}`)}
			<div class="grid rounded-lg border text-center">
				<p
					class="border-b-muted-foreground/30 bg-muted text-muted-foreground z-10 rounded-t-lg border-b px-6 py-3 font-semibold shadow-sm"
				>
					{hour.weekday}
				</p>
				<div class="grid gap-1">
					<Button
						type="button"
						variant="ghost"
						class="h-6 cursor-pointer rounded-none border-b px-1 py-1 has-[>svg]:p-1"
						onclick={() => buttonUpClicked(index)}
						disabled={$submitting}
					>
						<ChevronUp />
					</Button>
					<p class="py-1 text-sm leading-none">{toDurationString(hour.hours)}</p>
					<Button
						type="button"
						variant="ghost"
						class="h-6 cursor-pointer rounded-t-none border-t px-1 py-1 has-[>svg]:p-1"
						onclick={() => buttonDownClicked(index)}
						disabled={$submitting}
					>
						<ChevronDown class="size-4" />
					</Button>
				</div>
			</div>
		{/each}
	</div>
</form>
