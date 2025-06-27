<script lang="ts">
	import { superForm } from 'sveltekit-superforms';
	import type { PageProps } from './$types';
	import { Label } from '$lib/components/ui/label';
	import { Input } from '$lib/components/ui/input';
	import { Button } from '$lib/components/ui/button';

	let { data }: PageProps = $props();

	const { form } = superForm(data.form);
</script>

<form method="POST" class="flex w-[350px] flex-col gap-6">
	<div class="flex flex-col items-center gap-2 text-center">
		<h1 class="text-2xl font-bold">Opret din konto</h1>
		<p class="text-muted-foreground text-sm text-balance">
			Indtast din information her for at oprette dig
		</p>
	</div>
	<div class="grid gap-6">
		<div class="flex flex-col gap-2">
			<Label for="name">Navn</Label>
			<Input type="text" name="name" bind:value={$form.name} />
		</div>
		<div class="flex flex-col gap-2">
			<Label for="email">Email</Label>
			<Input type="email" name="email" bind:value={$form.email} />
		</div>
		<div class="flex flex-col gap-2">
			<Label for="password">Adgangskode</Label>
			<Input type="password" name="password" bind:value={$form.password} />
		</div>
		<Button type="submit">Registrer</Button>
	</div>
	<div class="text-center text-sm">
		Har du allerede en konto?
		<a href="/auth/login" class="underline underline-offset-4"> Log ind her </a>
	</div>
</form>
