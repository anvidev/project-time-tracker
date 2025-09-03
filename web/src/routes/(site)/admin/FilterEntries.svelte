<script lang="ts">
	import * as Accordion from '$lib/components/ui/accordion/index.js';
</script>

<Accordion.Root type="single" class="rounded-md border px-3 py-2">
	<Accordion.Item value="item-1">
		<Accordion.Trigger>Filtrer</Accordion.Trigger>
		<Accordion.Content>Yes. It adheres to the WAI-ARIA design pattern.</Accordion.Content>
	</Accordion.Item>
</Accordion.Root>
