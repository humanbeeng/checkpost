<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Avatar from '@/components/ui/avatar';
	import { Button } from '@/components/ui/button';
	import type { User } from '@/types';

	export let user: User;
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger asChild let:builder>
		<Button variant="ghost" builders={[builder]} class="relative rounded-md hover:bg-gray-500/10">
			<Avatar.Root class="h-8 w-8">
				<Avatar.Image src={user.avatar_url} alt={user.name} />
				<!-- <Avatar.Fallback></Avatar.Fallback> -->
			</Avatar.Root>
			<p class="px-2" autocapitalize="on">{user.name}</p>
		</Button>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content class="w-56" align="end">
		<a href="/auth/logout" data-sveltekit-reload>
			<DropdownMenu.Item>Log out</DropdownMenu.Item>
		</a>
	</DropdownMenu.Content>
</DropdownMenu.Root>
