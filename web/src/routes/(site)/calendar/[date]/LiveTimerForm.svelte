<script lang="ts">
	import { enhance } from '$app/forms';
	import { format } from 'date-fns';
	import {
		Select,
		SelectContent,
		SelectGroup,
		SelectItem,
		SelectLabel,
		SelectSeparator,
		SelectTrigger
	} from '$lib/components/ui/select';
	import { Label } from '$lib/components/ui/label';
	import { Button } from '$lib/components/ui/button';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { LoaderCircle, Pause, Play, Save, Trash2 } from '@lucide/svelte';
	import type { Category, LiveTimer } from '$lib/types';
	import { Millisecond } from '$lib/utils';

	const { timer, categories }: { timer: LiveTimer | null; categories: Category[] } = $props();

	let formEl: HTMLFormElement;
	let updateButton: HTMLButtonElement;

	let categoryId = $state(-1);
	let description = $state('');
	let busy = $state(false);
	let error = $state('');

	// keep the description in sync with the timer, which can also be changed from other clients
	$effect(() => {
		description = timer?.description ?? '';
	});

	// the server only tells how much time was tracked when the timer was loaded,
	// so the running time is counted from that moment on in the browser
	const base = $derived(timer ? { elapsed: timer.elapsed, at: Date.now() } : null);

	let now = $state(Date.now());
	$effect(() => {
		if (timer?.status != 'running') {
			return;
		}

		now = Date.now();
		const interval = setInterval(() => (now = Date.now()), 250);
		return () => clearInterval(interval);
	});

	const elapsedSeconds = $derived.by(() => {
		if (!base || !timer) {
			return 0;
		}

		const running = timer.status == 'running' ? Math.max(now - base.at, 0) : 0;
		return Math.floor((base.elapsed / Millisecond + running) / 1000);
	});

	const clock = $derived.by(() => {
		const pad = (n: number) => n.toString().padStart(2, '0');
		const hours = Math.floor(elapsedSeconds / 3600);
		const minutes = Math.floor((elapsedSeconds % 3600) / 60);

		return `${pad(hours)}:${pad(minutes)}:${pad(elapsedSeconds % 60)}`;
	});

	const categoryMap = $derived.by(() => {
		const res: Record<string, Category[]> = {};
		for (const category of categories) {
			(res[category.rootTitle] ??= []).push(category);
		}

		return res;
	});

	const categoryTitle = $derived.by(() => {
		const id = timer ? timer.categoryId : categoryId;
		return categories.find((c) => c.id == id)?.title ?? (timer ? `Kategori ${id}` : 'Vælg en kategori');
	});

	const submitEnhance: Parameters<typeof enhance>[1] = ({ formData }) => {
		// the date is decided by the browser, so it matches the local day of the user
		formData.set('date', format(new Date(), 'yyyy-MM-dd'));
		busy = true;

		return async ({ result, update }) => {
			busy = false;
			error =
				result.type == 'failure'
					? ((result.data?.error as string | undefined) ?? 'Noget gik galt')
					: '';

			await update({ reset: false });
		};
	};

	const saveDescription = () => {
		if (timer && description != timer.description) {
			formEl.requestSubmit(updateButton);
		}
	};
</script>

<form
	bind:this={formEl}
	method="POST"
	class="flex h-full flex-col justify-between gap-4"
	use:enhance={submitEnhance}
>
	<div class="flex flex-col gap-4">
		<div class="grid gap-1">
			<Label class="gap-[2px]" for="live-category">
				Kategori{#if !timer}<span class="text-red-700">*</span>{/if}
			</Label>
			<input type="hidden" name="categoryId" value={categoryId} />
			<Select
				type="single"
				name="live-category"
				disabled={timer != null}
				value={categoryId.toString()}
				onValueChange={(id) => (categoryId = parseInt(id))}
			>
				<SelectTrigger class="w-full" id="live-category">
					{categoryTitle}
				</SelectTrigger>
				<SelectContent>
					{#each Object.entries(categoryMap) as [label, items], index (label)}
						<SelectGroup>
							{#if index > 0}
								<SelectSeparator class="ml-px" />
							{/if}
							<SelectLabel>{label}</SelectLabel>
							{#each items as category (category.id)}
								<SelectItem class="pl-4" value={category.id.toString()}>
									{category.title}
								</SelectItem>
							{/each}
						</SelectGroup>
					{/each}
				</SelectContent>
			</Select>
		</div>

		<div class="grid gap-1">
			<Label for="live-description">Beskrivelse</Label>
			<Textarea
				id="live-description"
				name="description"
				bind:value={description}
				onblur={saveDescription}
				placeholder="Indtast valgfri beskrivelse"
				class="resize-none"
			></Textarea>
		</div>
	</div>

	<div class="flex flex-col items-center gap-3">
		<div
			class="text-4xl font-semibold tabular-nums {timer?.status != 'running'
				? 'text-muted-foreground'
				: ''}"
			aria-live="off"
		>
			{clock}
		</div>

		<!-- the same controls are shown before and after the timer is started, so nothing jumps around -->
		<div class="grid w-full grid-cols-3 gap-2">
			{#if timer?.status == 'running'}
				<Button
					type="submit"
					formaction="?/pauseTimer"
					class="cursor-pointer"
					disabled={busy}
				>
					<Pause /> Pause
				</Button>
			{:else if timer}
				<Button
					type="submit"
					formaction="?/resumeTimer"
					class="cursor-pointer"
					disabled={busy}
				>
					<Play /> Genoptag
				</Button>
			{:else}
				<Button
					type="submit"
					formaction="?/startTimer"
					class="cursor-pointer"
					disabled={busy || categoryId < 0}
				>
					<Play /> Start
				</Button>
			{/if}
			<Button
				type="submit"
				formaction="?/saveTimer"
				variant="outline"
				class="cursor-pointer"
				disabled={busy || !timer}
			>
				<Save /> Gem
			</Button>
			<Button
				type="submit"
				formaction="?/discardTimer"
				variant="ghost"
				class="cursor-pointer"
				disabled={busy || !timer}
				onclick={(e: MouseEvent) => {
					if (!confirm('Kassér timeren uden at gemme?')) {
						e.preventDefault();
					}
				}}
			>
				<Trash2 /> Kassér
			</Button>
		</div>

		<p class="text-muted-foreground text-center text-xs">
			Når du gemmer, bliver tiden registreret, og timeren kan ikke genoptages.
		</p>
		<p class="min-h-5 text-center text-sm text-red-700">
			{#if busy}
				<LoaderCircle class="inline animate-spin" />
			{:else}
				{error}
			{/if}
		</p>
	</div>

	<button
		type="submit"
		formaction="?/updateTimer"
		aria-label="Gem beskrivelse"
		hidden
		bind:this={updateButton}
	></button>
</form>
