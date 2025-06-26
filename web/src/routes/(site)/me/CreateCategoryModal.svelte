<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { cn } from '$lib/utils';
	import { Loader2, Plus } from '@lucide/svelte';
	import { superForm } from 'sveltekit-superforms';

	const {
		triggerClass,
		parentId,
		parentName,
		disabled = false
	}: {
		triggerClass: string;
		parentId: number | null;
		parentName: string | null;
		disabled?: boolean;
	} = $props();

	let open = $state(false);

	const { form, formId, submitting, delayed } = superForm(
		{
			parentId,
			title: ''
		},
		{
			onUpdated: () => {
				open = false;
			},
			delayMs: 350,
			timeoutMs: 8000,
			multipleSubmits: 'prevent'
		}
	);
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger {disabled} class={cn(triggerClass, 'bg-muted/50 hover:bg-muted cursor-pointer')}>
		<Plus class="size-4" />
	</Dialog.Trigger>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Opret en ny kategori</Dialog.Title>
			{#if parentName != null}
				<Dialog.Description>Kategorien er en underkategori til {parentName}</Dialog.Description>
			{/if}
		</Dialog.Header>
		<form method="POST" action="?/createCategory" id={$formId} class="flex flex-col gap-4">
			<div class="grid gap-2">
				<Label class="gap-[2px]" for="title">Titel<span class="text-red-700">*</span></Label>
				<Input
					class="input-arrows-none"
					name="title"
					bind:value={$form.title}
					type="text"
					placeholder="Indtast kategoriens titel"
				/>
			</div>
			<Input type="hidden" name="parentId" value={$form.parentId} />
		</form>
		<Dialog.Footer>
			<Button disabled={$submitting} type="submit" form={$formId}>
				{#if $delayed}
					<Loader2 class="animate-spin" />
				{:else}
					<Plus class="size-4" />
				{/if}
				Opret
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
