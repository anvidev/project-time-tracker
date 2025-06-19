<script lang="ts">
	import type { LayoutProps } from './$types';

	const { children }: LayoutProps = $props();
</script>

<main
	class="mx-auto flex min-h-[100vh] max-w-[1600px] flex-col items-center justify-center overflow-auto px-8"
>
	{@render children()}
</main>
