<script lang="ts">
	import * as Accordion from '$lib/components/ui/accordion/index.js';
	import type { Category, User } from '$lib/types';
	import { queryParamsState } from 'kit-query-params';

	const queryParams = queryParamsState({
		schema: {
			query: 'string',
			categoryId: ['number'],
			userId: ['number'],
			fromDate: 'date',
			toDate: 'date'
		}
	});

	const { users, categories }: { users: User[]; categories: Category[] } = $props();
</script>

<pre>
    {JSON.stringify(queryParams, null, 2)}
</pre>

<Accordion.Root type="single" class="rounded-md border px-3 py-2">
	<Accordion.Item value="item-1">
		<Accordion.Trigger>Filtrer</Accordion.Trigger>
		<Accordion.Content>
			<input bind:value={queryParams.query} />
		</Accordion.Content>
	</Accordion.Item>
</Accordion.Root>
