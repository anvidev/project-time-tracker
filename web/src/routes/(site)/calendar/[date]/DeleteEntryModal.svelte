<script lang="ts">
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button';
	import { LoaderCircle, Trash } from '@lucide/svelte';
	import { superForm } from 'sveltekit-superforms';

	const data: { id: number } = $props();

	const { form, enhance, delayed } = superForm(data, {
		delayMs: 150,
		timeoutMs: 8000
	});
</script>

<AlertDialog.Root>
	<AlertDialog.Trigger
		class={buttonVariants({
			variant: 'outline',
			class: 'cursor-pointer text-red-500 hover:text-red-500'
		})}
	>
		<Trash />
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Er du sikker?</AlertDialog.Title>
			<AlertDialog.Description>
				Denne registrering vil blive slettet for evigt og kan ikke gendannes. Vil du fortsætte?
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Annuller</AlertDialog.Cancel>
			<form method="POST" action="?/deleteTimeEntry" use:enhance>
				<input type="hidden" name="id" bind:value={$form.id} />
				<AlertDialog.Action type="submit" class={buttonVariants({ variant: 'destructive', class: "cursor-pointer" })}>
					Slet
					{#if $delayed}
						<LoaderCircle class="animate-spin" />
					{:else}
						<Trash />
					{/if}
				</AlertDialog.Action>
			</form>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
