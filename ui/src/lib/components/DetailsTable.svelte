<script lang="ts">
	export let data: Object;
	export let title: string;
</script>

<div class="border py-2 px-4 rounded bg-gray-100 shadow-sm my-2">
	<h4 class="font-medium text-md">{title}</h4>
	{#if data}
		{#each Object.entries(data) as d}
			<hr />
			<div class="grid grid-cols-5 text-sm w-full my-1 gap-2">
				<p class="col-span-1 text-wrap">{d[0]}</p>
				<code
					class="col-span-4 overflow-hidden hover:overflow-auto text-ellipsis whitespace-normal"
				>
					{d[1]}
				</code>
			</div>
		{/each}
	{:else}
		<hr />
		<p>(empty)</p>
	{/if}
</div>
