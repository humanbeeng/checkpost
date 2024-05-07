<script lang="ts">
	import Header from '@/components/Header.svelte';

	import { enhance } from '$app/forms';
	import { PUBLIC_SERVER_URL } from '$env/static/public';
	import { Button } from '@/components/ui/button';
	import { debounce } from '@/debounce';
	import { Link1, LinkBreak2, Reload } from 'svelte-radix';
	import type { EndpointExistsResponse, State } from '../types';

	export let form;
	export let data;
	let submitBtnText = 'Continue';

	let subdomainError: string | null;
	let endpointExistsResponse: EndpointExistsResponse | null;
	let subdomain: string = '';

	let state: State = 'empty';

	async function checkSubdomain(event: any) {
		subdomain = event.target.value;
		if (!subdomain) {
			state = 'empty';
			return;
		}

		// Call check api
		if (subdomain.toString().length < 4 || subdomain.toString().length > 10) {
			state = 'error';
			return;
		}

		const res = await fetch(`${PUBLIC_SERVER_URL}/url/exists/${subdomain}`);

		switch (res.status) {
			case 200: {
				endpointExistsResponse = (await res.json()) as EndpointExistsResponse;
				subdomainError = null;
				state = 'success';
				break;
			}
			case 400: {
				subdomainError = await res.text();
				state = 'error';
				endpointExistsResponse = null;
				break;
			}
		}
	}

	const debouncedHandleChange = debounce(checkSubdomain, 1000);

	function handleInput(e: any) {
		if (e.which === 32) {
			// 32 is the keycode for space
			e.preventDefault();
			return;
		}
		if ((subdomain = e.target.value)) {
			state = 'loading';
			submitBtnText = 'Continue';
		}
		debouncedHandleChange(e);
	}
</script>

<body class="h-screen flex flex-col bg-gray-40">
	<Header user={data.user ?? null} />

	<main class="w-full items-center flex flex-col flex-grow justify-center bg-gray-100/10">
		<div
			class="border py-3 px-4 rounded-xl mb-52 flex flex-col items-center justify-center gap-6 bg-white/90 shadow-md hover:isolate"
		>
			<div class="w-full">
				<h3 class="self-start font-medium text-lg px-1">Confirm your endpoint</h3>
				<hr class="mt-1" />
			</div>
			<div class="flex gap-2 p-1">
				<form
					method="POST"
					use:enhance={() => {
						state = 'loading';
						submitBtnText = 'Confirming';
						return async ({ update }) => {
							await update();
							state = 'empty';
						};
					}}
				>
					<div class="flex flex-col gap-2">
						<div class="flex flex-col">
							<span class="flex gap-1">
								{#if state === 'error' || endpointExistsResponse?.exists}
									<LinkBreak2 class="w-8" />
								{:else}
									<Link1 class="w-8" />
								{/if}
								<p>https://</p>
								<input
									bind:value={subdomain}
									spellcheck="false"
									autocomplete="off"
									id="subdomain"
									name="subdomain"
									placeholder="dunder-mifflin"
									type="text"
									class="border-b border-b-green outline-none w-28"
									maxlength="10"
									minlength="4"
									on:keydown={handleInput}
									pattern="[A-Za-z0-9]+"
								/>.checkpost.io
							</span>

							<span class=" mb-2 h-4 self-end">
								{#if state === 'success' && endpointExistsResponse}
									{#if !endpointExistsResponse.exists}
										<p class="text-green-900 text-xs py-1">It's available</p>
									{:else}
										<p class="text-red-900 text-xs py-1">{endpointExistsResponse?.message}</p>
									{/if}
								{/if}
								{#if state === 'error'}
									<p class="text-red-900 text-xs py-2">{subdomainError}</p>
								{/if}
							</span>
						</div>
						<!-- <label for="subdomain">Your URL</label> -->
						<Button class="w-full" disabled={state !== 'success'} variant="default" type="submit">
							{#if state === 'loading'}
								<Reload class="mr-2 h-4 w-4 animate-spin" />
							{/if}
							{submitBtnText}
						</Button>

						{#if form?.err}<p class="text-red-950 text-xs max-w-72">{form.err.message}</p>{/if}
					</div>
				</form>
			</div>
		</div>
	</main>
</body>
