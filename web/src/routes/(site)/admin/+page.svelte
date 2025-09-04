<script lang="ts">
	import PieChartEntriesCount from './PieChartEntriesCount.svelte';
	import PieChartEntriesHours from './PieChartEntriesHours.svelte';
	import PieChartEntriesUsers from './PieChartEntriesUsers.svelte';
	import BarChartEntriesCount from './BarChartEntriesCount.svelte';
	import TableEntries from './TableEntries.svelte';
	import FilterEntries from './FilterEntries.svelte';

	const { data } = $props();
	const { entries, categories, users } = $derived(data);
</script>

<div class="w-full">
	<FilterEntries {users} {categories} />
</div>

<div class="flex w-full flex-col gap-4 lg:flex-row lg:items-center">
	<div class="w-1/3">
		<PieChartEntriesCount entries={entries.entries} />
	</div>

	<div class="w-1/3">
		<PieChartEntriesHours entries={entries.entries} timeSpent={entries.timeSpent} />
	</div>

	<div class="w-1/3">
		<PieChartEntriesUsers entries={entries.entries} />
	</div>
</div>

<div class="w-full">
	<BarChartEntriesCount entries={entries.entries} />
</div>

<TableEntries entries={entries.entries} />
