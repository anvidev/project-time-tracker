<script lang="ts">
	import type { LayoutProps } from './$types';

	const { children }: LayoutProps = $props();
</script>

<main
	class="mx-auto flex min-h-screen max-w-7xl gap-4 flex-col items-center overflow-auto py-6 px-8"
>
	{@render children()}
</main>
