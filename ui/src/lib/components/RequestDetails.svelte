<script lang="ts">
	import { isFormatEnabled } from '@/store';
	import type { Request } from '@/types';
	import { extraHeaders, formatJson, timeAgo } from '@/utils.js';
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

<div class="grid grid-cols-1 w-full h-fit gap-6">
	<!-- Request details -->
	<div class="border py-2 px-3 rounded bg-gray-100 shadow-sm my-2">
		<h4 class="font-medium text-md">Request details</h4>
		<hr />
		<div class="grid grid-cols-10 gap-3 my-1 text-sm">
			<p>ID</p>
			<code class="overflow-hidden col-span-9 font-sans">{request.uuid}</code>
		</div>

		<hr />
		<div class="grid grid-cols-10 gap-3 my-1 text-sm">
			<p>Path</p>
			<code class="overflow-hidden col-span-9 font-sans">{request.path}</code>
		</div>

		<hr />
		<div class="grid grid-cols-10 gap-3 my-1 text-sm">
			<p>Method</p>
			<code class="font-sans overflow-hidden col-span-9">{request.method.toUpperCase()}</code>
		</div>

		<hr />
		<div class="grid grid-cols-10 gap-3 my-1 text-sm">
			<p>Size</p>
			<code class="overflow-hidden col-span-9 font-sans">{request.content_size}</code>
		</div>

		<hr />
		<div class="grid grid-cols-10 gap-3 my-1 text-sm">
			<p>Time</p>
			<code class="overflow-hidden col-span-9 font-sans"
				>{request.created_at} ({timeAgo(request.created_at)})</code
			>
		</div>
	</div>

	<!-- Headers-->
	<DetailsTable
		title="Headers"
		data={request.headers}
		showHideButton={true}
		hiddenEntries={extraHeaders}
	/>

	{#if request.query_params}
		<div class="grid grid-cols-1 lg:grid-cols-2 w-full gap-3">
			<!-- Query params-->
			<DetailsTable title="Query" data={request.query_params} />

			<!-- Form values-->
			<DetailsTable title="Form" data={request.form_data} />
		</div>
	{/if}
</div>
<!-- File attachments -->
<!-- TODO: Add file attachments -->

<!-- <DetailsTable title="Form" data={request.query_params} /> -->
<hr class="mt-4 mb-2" />

<!-- Payload -->
<div class="mb-6 border py-2 px-4 rounded-md bg-gray-100">
	<div class="flex justify-between mb-1">
		<h4 class="font-medium text-md">Payload</h4>

		{#if request.content && !(request.content_type.startsWith('multipart/form-data') || request.content_type.startsWith('application/x-www-form-urlencoded'))}
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

	{#if content && !(request.content_type.startsWith('multipart/form-data') || request.content_type.startsWith('application/x-www-form-urlencoded'))}
		{#if checked}
			<pre class="bg-gray-50 border rounded p-4 shadow-sm overflow-x-scroll"><code
					class="font-mono text-sm">{prettyContent}</code
				>
		</pre>
		{:else}
			<pre class="bg-gray-50 border rounded p-4 shadow-sm overflow-x-scroll"><code
					class="font-mono text-sm">{content}</code
				>
		</pre>
		{/if}
	{:else}
		<hr class="my-2" />

		<div class="w-full flex flex-col place-items-center align-middle justify-center min-h-6">
			<p class="font-light text-xs py-2">(empty)</p>
		</div>
	{/if}
</div>
