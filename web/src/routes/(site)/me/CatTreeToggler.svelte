<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Circle, CircleCheck, CircleDot, Loader2 } from '@lucide/svelte';
	import { superForm } from 'sveltekit-superforms';

	let {
		parentIsFollowed,
		isNestedChildFollowed,
		onSubmittingChange,
		submitting = $bindable(false),
		...formData
	}: {
		id: number;
		isFollowed: boolean;
		parentIsFollowed: boolean;
		isNestedChildFollowed: boolean;
		submitting?: boolean;
		onSubmittingChange?: (v: boolean) => void;
	} = $props();

	const { form, enhance, submitting: formSubmitting, delayed } = superForm(formData, {
		delayMs: 250,
		timeoutMs: 8000,
		onUpdated: ({ form: updatedForm }) => {
			$form.isFollowed = !updatedForm.data.isFollowed
		},
	});

	let lastSubmittingValue = $state($formSubmitting)

	$effect(() => {
		if (lastSubmittingValue != $formSubmitting) {
			if(submitting != undefined) 
				submitting = $formSubmitting
				
			onSubmittingChange?.($formSubmitting)
			lastSubmittingValue = $formSubmitting
		}
	})

	const { isFollowed } = $derived(formData)
</script>

<form method="POST" action="?/toggleFollow" use:enhance class="size-4">
	<input type="hidden" name="id" value={$form.id} />
	<input type="hidden" name="isFollowed" value={$form.isFollowed} />
	<Button
		disabled={$formSubmitting || parentIsFollowed}
		type="submit"
		variant="ghost"
		size="sm"
		class="size-fit cursor-pointer p-0 has-[>svg]:px-0 group"
	>
		{#if $delayed}
			<Loader2 class="size-4 animate-spin" />
		{:else if parentIsFollowed || isFollowed}
			<CircleCheck class="size-4 text-green-600 group-hover:!text-amber-700" />
		{:else if isNestedChildFollowed}
			<CircleDot class="size-4 text-yellow-600 group-hover:!text-amber-700" />
		{:else}
			<Circle class="size-4 group-hover:!text-amber-700" />
		{/if}
	</Button>
</form>
