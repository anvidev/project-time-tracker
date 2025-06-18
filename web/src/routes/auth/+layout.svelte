<script lang="ts">
	import type { LayoutProps } from './$types';

	const { children }: LayoutProps = $props();
</script>

<main class="flex h-[100vh] max-w-[1600px] flex-col items-center justify-center overflow-auto mx-auto">
	{@render children()}
</main>
