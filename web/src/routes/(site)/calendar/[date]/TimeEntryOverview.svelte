<script lang="ts">
	import type { SummaryDay } from '$lib/types';
	import { Hour } from '$lib/utils';
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

		return daySummary.totalHours / Hour;
	});

	const maxHours = $derived.by(() => {
		if (daySummary.maxHours == 0) {
			return 0;
		}

		return daySummary.maxHours / Hour;
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
							{((entry.duration / daySummary.totalHours) * 100).toFixed(1)}% af total
						</p>
					</div>
					<div class="flex gap-2">
						<Badge class="bg-background" variant="outline">
							{(entry.duration / Hour).toFixed(2)}t
						</Badge>
						<DeleteEntryModal id={entry.id} />
					</div>
				</div>
			{/each}
		</CardContent>
	</Card>
{/if}
