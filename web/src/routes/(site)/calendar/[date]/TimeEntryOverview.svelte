<script lang="ts">
	import type { SummaryDay } from '$lib/types';
	import { Hour, maxFractionDigits } from '$lib/utils';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import DeleteEntryModal from './DeleteEntryModal.svelte';

	const {
		daySummary,
		formattedDate
	}: {
		daySummary: SummaryDay;
		formattedDate: string;
	} = $props();

	const totalHours = $derived.by(() => {
		if (daySummary.totalHours == 0) {
			return 0;
		}

		return maxFractionDigits(daySummary.totalHours / Hour, 2);
	});

	const maxHours = $derived.by(() => {
		if (daySummary.maxHours == 0) {
			return 0;
		}

		return maxFractionDigits(daySummary.maxHours / Hour, 2);
	});
</script>

{#if daySummary.timeEntries.length > 0}
	<Card class="col-span-2">
		<CardHeader>
			<CardTitle>
				Registreringer d. {formattedDate}
			</CardTitle>
			<CardDescription>
				{totalHours} af {maxHours} timer registreret
			</CardDescription>
		</CardHeader>
		<CardContent>
			{#each daySummary.timeEntries as entry (entry.id)}
				<div class="bg-muted mb-2 flex items-center justify-between rounded-lg border p-2">
					<div>
						<p class="font-semibold tracking-tight">{entry.category}</p>
						<p class="text-muted-foreground text-sm">
							{maxFractionDigits((entry.duration / daySummary.totalHours) * 100, 2)}% af total
						</p>
					</div>
					<div class="flex gap-2">
						<Badge class="bg-background px-3 py-2" variant="outline">
							{(entry.duration / Hour).toFixed(2)}t
						</Badge>
						<DeleteEntryModal id={entry.id} />
					</div>
				</div>
			{/each}
		</CardContent>
	</Card>
{/if}
