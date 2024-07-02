<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Dot, DotFilled } from 'svelte-radix';

	export let state: 'connecting' | 'online' | 'offline' = 'connecting';
</script>

{#if state === 'online'}
	<Badge
		variant="outline"
		class="bg-green-100 text-green-600 dark:bg-green-900/20 dark:text-green-400 py-1 max-h-6 px-0"
		><span class="flex place-items-center">
			<DotFilled />
			<p class="pr-2">Listening</p>
		</span></Badge
	>
{:else if state === 'offline'}
	<Badge
		variant="outline"
		class="bg-red-100 text-red-600 py-1 max-h-6 dark:bg-red-900/20  dark:text-red-400 px-0"
	>
		<span class="flex place-items-center">
			<Dot />
			<p class="pr-2">Offline</p>
		</span>
	</Badge>
{:else if state === 'connecting'}
	<Badge
		variant="outline"
		class="bg-red-100 text-red-600 py-1 max-h-6 dark:bg-red-900/20  dark:text-red-400 px-0"
	>
		<span class="flex place-items-center">
			<Dot />
			<p class="pr-2">Connecting...</p>
		</span>
	</Badge>
{/if}
