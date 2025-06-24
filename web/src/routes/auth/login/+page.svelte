<script lang="ts">
	import { superForm } from 'sveltekit-superforms';
	import type { PageProps } from './$types';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';

	let { data }: PageProps = $props();

	const { form } = superForm(data.form);
</script>

<form method="POST" class="flex w-[350px] flex-col gap-3 rounded border p-4 shadow-sm">
	<div class="flex flex-col gap-1">
		<Label for="email">Email</Label>
		<Input type="email" name="email" bind:value={$form.email} />
	</div>
	<div class="flex flex-col gap-1">
		<Label for="password">Adgangskode</Label>
		<Input type="password" name="password" bind:value={$form.password} />
	</div>
	<Button type="submit" class="mt-2">Log ind</Button>
</form>
