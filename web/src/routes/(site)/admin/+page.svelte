<script lang="ts">
	import PieChartEntriesCount from './PieChartEntriesCount.svelte';
	import PieChartEntriesHours from './PieChartEntriesHours.svelte';
	import PieChartEntriesUsers from './PieChartEntriesUsers.svelte';
	import BarChartEntriesCount from './BarChartEntriesCount.svelte';
	import TableEntries from './TableEntries.svelte';
	import FilterEntries from './FilterEntries.svelte';
	import TableCategories from './TableCategories.svelte';

	const { data } = $props();
	const { entries, categories, users } = $derived(data);
</script>

<div class="w-full">
	<FilterEntries {users} {categories} />
</div>

<div class="flex w-full flex-col gap-4 lg:flex-row lg:items-center">
	<div class="lg:w-1/3">
		<PieChartEntriesCount entries={entries.entries} />
	</div>

	<div class="lg:w-1/3">
		<PieChartEntriesHours entries={entries.entries} timeSpent={entries.timeSpent} />
	</div>

	<div class="lg:w-1/3">
		<PieChartEntriesUsers entries={entries.entries} />
	</div>
</div>

<div class="w-full">
	<BarChartEntriesCount entries={entries.entries} />
</div>

<div class="mt-2 w-full space-y-3">
	<h2 class="font-semibold leading-none">Timer fordelt på kategorier</h2>
	<TableCategories entries={entries.entries} />
</div>

<div class="mt-2 w-full space-y-3">
	<h2 class="font-semibold leading-none">Tidsregistreringer</h2>
	<TableEntries entries={entries.entries} />
</div>
