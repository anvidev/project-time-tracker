<script lang="ts">
	import Input from '$lib/components/ui/input/input.svelte';
	import Label from '$lib/components/ui/label/label.svelte';
	import type { Category, User } from '$lib/types';
	import { queryParamsState } from 'kit-query-params';
	import { Button } from '$lib/components/ui/button';
	import FilterCategories from './FilterCategories.svelte';
	import FilterUsers from './FilterUsers.svelte';
	import FilterDate from './FilterDate.svelte';

	const queryParams = queryParamsState({
		schema: {
			query: 'string',
			categoryId: ['number'],
			userId: ['number'],
			fromDate: 'date',
			toDate: 'date'
		},
		debounce: 500
	});

	const { users, categories }: { users: User[]; categories: Category[] } = $props();
</script>

<div class="flex items-center gap-2">
	<div class="grid gap-1.5">
		<Label>Beskrivelse</Label>
		<Input bind:value={queryParams.query} placeholder="Søg i beskrivelse..." />
	</div>
	<FilterCategories {categories} bind:param={queryParams.categoryId} />
	<FilterUsers {users} bind:param={queryParams.userId} />
	<FilterDate bind:param={queryParams.fromDate} label="Fra dato" />
	<FilterDate bind:param={queryParams.toDate} label="Til dato" />
	<Button variant="ghost" class="self-end" onclick={() => queryParams.$reset()}>Nulstil</Button>
</div>
