<script lang="ts">
	import type { LayoutProps } from './$types';

	const { children }: LayoutProps = $props();
</script>

<main
	class="mx-auto flex min-h-[100vh] max-w-[1600px] flex-col items-center overflow-auto py-6 px-8"
>
	{@render children()}
</main>
