<script lang="ts">
	import { buttonVariants } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { User } from '@lucide/svelte';
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class={buttonVariants({ variant: 'outline', size: 'icon', class: 'ml-auto cursor-pointer' })}
	>
		<User />
	</DropdownMenu.Trigger>
	<DropdownMenu.Content>
		<DropdownMenu.Item class="cursor-pointer">
			{#snippet child({props})}
				<a {...props} href="/me">Profil</a>
			{/snippet}
		</DropdownMenu.Item>
		<DropdownMenu.Separator />
		<DropdownMenu.Item class="hover:!text-destructive transition-colors cursor-pointer">
			{#snippet child({props})}
				<a {...props} href="/auth/logout">Log ud</a>
			{/snippet}
		</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root>
