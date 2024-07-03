<script lang="ts">
	import { removeKeys } from '@/utils.js';
	import Button from './ui/button/button.svelte';

	export let data: Object;
	export let title: string;
	export let showHideButton = false;
	export let hiddenEntries: string[] = [];
	let showHiddenEntries = false;

	$: displayedData = data;
	$: {
		const filteredData = removeKeys(data, hiddenEntries);
		if (showHiddenEntries) {
			displayedData = data;
		} else {
			displayedData = filteredData;
		}
	}
</script>

<div class="border py-2 px-4 rounded bg-gray-100 shadow-sm my-2">
	<span class="flex align-middle justify-between w-full py-1">
		<h4 class="font-medium text-md">{title}</h4>
		{#if showHideButton}
			<Button
				class="h-6 text-xs font-light px-1 min-w-20"
				variant="outline"
				on:click={() => {
					showHiddenEntries = !showHiddenEntries;
				}}
			>
				{#if showHiddenEntries}
					Hide
				{:else}
					Show hidden
				{/if}
			</Button>
		{/if}
	</span>
	{#if displayedData && Object.keys(displayedData).length}
		{#each Object.entries(displayedData) as d}
			<hr />
			<div class="grid grid-cols-5 text-sm w-full my-1 gap-2">
				<p class="col-span-1 text-wrap">{d[0]}</p>
				<code class="col-span-4 overflow-hidden hover:overflow-auto whitespace-normal font-sans">
					{d[1]}
				</code>
			</div>
		{/each}
	{:else}
		<hr />
		<div class="w-full flex flex-col place-items-center align-middle justify-center min-h-6">
			<p class="font-light text-xs py-2">(empty)</p>
		</div>
	{/if}
</div>
