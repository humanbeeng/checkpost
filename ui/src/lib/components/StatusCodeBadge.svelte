<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index.js';

	export let response_code: number;
</script>

{#if response_code >= 200 && response_code < 300}
	<Badge
		variant="outline"
		class="bg-green-100 text-green-600 dark:bg-green-900/20 dark:text-green-400 py-1 max-h-6"
	>
		<span class="flex gap-1">
			<p>
				{response_code}
			</p>
			<p>OK</p>
		</span>
	</Badge>
{:else}
	<Badge variant="outline" class="bg-red-100 text-red-600 dark:bg-red-900/20  dark:text-red-400"
		>{response_code}</Badge
	>
{/if}
