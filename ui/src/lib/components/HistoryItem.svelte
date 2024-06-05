<script lang="ts">
	import type { Request } from '@/types';
	import clsx from 'clsx';
	import MethodBadge from './MethodBadge.svelte';

	export let request: Request;
</script>

<div class="border-b border-gray-300 border-x-0 py-2 flex flex-col gap-1 w-full">
	<div class="flex justify-between">
		<MethodBadge method={request.method} />
		<p
			class={clsx(
				request.response_code === 200 && 'text-green-800',
				request.response_code !== 200 && 'text-red-800'
			)}
		>
			{request.response_code}
		</p>
	</div>
	<p class="opacity-70 text-start truncate">
		{#if request.path === '/'}
			{request.path}
		{:else}
			/{request.path}
		{/if}
	</p>
</div>
