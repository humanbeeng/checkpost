<script lang="ts">
	import Header from '@/components/Header.svelte';
	import { Check, Copy } from 'svelte-radix';

	export let data;
	export let url: string = 'https://google.checkpost.io/';
	let copied = false;

	const copy = () => {
		navigator.clipboard
			.writeText(url)
			.then(() => {
				copied = true;
			})
			.catch((err: any) => {
				console.error('Failed to copy text: ', err);
			});
	};
</script>

<body class="h-screen flex flex-col bg-gray-50">
	<Header user={data.user ?? null} />
	<div class="flex flex-col items-center justify-center w-full h-full gap-4 align-middle mb-36">
		<div class=" shadow-sm bg-zinc-200/70 rounded-sm border-lg flex py-2 px-3 gap-2">
			<button on:click={copy} class="hover:shadow-lg rounded-lg">
				{#if !copied}
					<Copy class="outline-none border-none p-1" />
				{:else}
					<Check class="outline-none" />
				{/if}
			</button>
			<p class="justify-center self-center">
				{url}
			</p>
		</div>
		<div class="flex flex-col items-center gap-4">
			<h1 class="font-medium tracking-tight text-2xl mt-2">Yay! Its yours🥳🎉</h1>
			<div class="flex flex-col items-center">
				<p>
					We are launching very soon! Stay tuned - we'll be notifying via mail you the moment it
					goes live.
				</p>
				<p>
					Meanwhile, catch us on
					<a
						href="https://twitter.com/checkposthq"
						target="_blank"
						rel="noreferrer"
						class="underline underline-offset-4"
						>X (@CheckpostHQ)
					</a> for the latest updates.
				</p>
			</div>
		</div>
	</div>
</body>
