<script lang="ts">
	import { isFormatEnabled } from '@/store';
	import type { Request } from '@/types';
	import { formatJson, timeAgo } from '@/utils.js';
	import DetailsTable from './DetailsTable.svelte';
	import { Checkbox } from './ui/checkbox';

	export let request: Request;
	let copied = false;

	$: content = request.content;
	$: prettyContent = formatJson(content);
	$: checked = $isFormatEnabled;

	const copy = (content: string) => {
		navigator.clipboard
			.writeText(content)
			.then(() => {
				copied = true;
				return true;
			})
			.catch((err: any) => {
				console.error('Failed to copy text: ', err);
			});
	};
</script>

<div class="grid grid-cols-1 w-full h-full gap-3">
	<div class="border py-2 px-3 rounded bg-gray-100 shadow-sm my-2">
		<h4 class="font-medium text-md">Request details</h4>
		<hr />
		<div class="grid grid-cols-10 gap-3 my-1 text-sm">
			<p>ID</p>
			<code class="overflow-hidden col-span-9">{request.uuid}</code>
		</div>

		<hr />
		<div class="grid grid-cols-10 gap-3 my-1 text-sm">
			<p>Path</p>
			<code class="overflow-hidden col-span-9">{request.path}</code>
		</div>

		<hr />
		<div class="grid grid-cols-10 gap-3 my-1 text-sm">
			<p>Method</p>
			<code class="overflow-hidden col-span-9">{request.method.toUpperCase()}</code>
		</div>

		<hr />
		<div class="grid grid-cols-10 gap-3 my-1 text-sm">
			<p>Size</p>
			<code class="overflow-hidden col-span-9">{request.content_size}</code>
		</div>

		<hr />
		<div class="grid grid-cols-10 gap-3 my-1 text-sm">
			<p>Time</p>
			<code class="overflow-hidden col-span-9"
				>{request.created_at} ({timeAgo(request.created_at)})</code
			>
		</div>
	</div>

	<!-- Headers-->
	<DetailsTable title="Headers" data={request.headers} />
</div>

{#if request.query_params}
	<div class="grid grid-col-1 lg:grid-cols-2 w-full gap-3">
		<!-- Query params-->

		<DetailsTable title="Query" data={request.query_params} />

		<!-- TODO: Add form values -->
		<!-- Form values-->
		<DetailsTable title="Form" data={request.query_params} />
	</div>
{/if}

<!-- File attachments -->
<!-- TODO: Add file attachments -->
<!-- <DetailsTable title="Form" data={request.query_params} /> -->
<hr class="mt-4 mb-2" />

<!-- Payload -->
<div class="mb-6 border py-2 px-4 rounded-md bg-gray-100">
	<div class="flex justify-between mb-1">
		<h4 class="font-medium text-md">Payload</h4>

		{#if request.content}
			<div class="inline-flex justify-center items-center gap-2">
				<span class="inline-flex place-items-center gap-1">
					<Checkbox id="terms" bind:checked class="border-gray-500" />
					<label
						for="terms"
						class="text-sm font-light leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
					>
						Format
					</label>
				</span>

				<button on:click={() => copy(content)} class="bg-gray-100 w-fit px-2 py-0">
					{#if !copied}
						<p class="font-light underline text-sm leading-none">Copy</p>
					{:else}
						<p class="font-light underline text-sm leading-none">Copied</p>
					{/if}
				</button>
			</div>
		{/if}
	</div>

	{#if content}
		{#if checked}
			<pre class="bg-gray-50 border rounded p-4 shadow-sm"><code>{prettyContent}</code>
		</pre>
		{:else}
			<pre class="bg-gray-50 border rounded p-4 shadow-sm"><code>{content}</code>
		</pre>
		{/if}
	{:else}
		<hr class="my-2" />
		<div class="flex w-full justify-center">
			<pre class="my-2">(empty)</pre>
		</div>
	{/if}
</div>