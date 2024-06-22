<script lang="ts">
	import { goto } from '$app/navigation';
	import background from '$lib/assets/background.svg';
	import { Button } from '$lib/components/ui/button/index.js';
	import Header from '@/components/Header.svelte';

	import { PUBLIC_BASE_URL } from '$env/static/public';
	import Footer from '@/components/Footer.svelte';
	import { debounce } from '@/debounce';
	import { ChevronRight, GithubLogo, Reload } from 'svelte-radix';
	import type { EndpointExistsResponse, State } from './types';

	let error: string | null;
	let existsRes: EndpointExistsResponse | null;
	let subdomain = '';

	let state: State = 'empty';

	async function checkSubdomain(event: any) {
		subdomain = event.target.value;
		if (!subdomain) {
			state = 'empty';
			return;
		}

		// Call check api
		if (subdomain.length < 4 || subdomain.length > 10) {
			error = 'Subdomain length should be between 4 and 10';
			state = 'error';
			return;
		}

		const res = await fetch(`${PUBLIC_BASE_URL}/endpoint/exists/${subdomain}`);

		switch (res.status) {
			case 200: {
				existsRes = (await res.json()) as EndpointExistsResponse;
				error = null;
				state = 'success';
				break;
			}
			case 400: {
				error = await res.text();
				state = 'error';
				existsRes = null;
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
		subdomain = e.target.value;
		if (subdomain.length < 4) {
			state = 'error';
		}
		state = 'loading';
		debouncedHandleChange(e);
	}

	const signin = () => {
		goto('/auth/github');
	};
</script>

<div class="w-full">
	<Header user={null}>
		<a
			href="https://checkpost.notion.site/Pricing-Proposed-c9399d8349f44ff9bf095817ad05bed0?pvs=74"
			target="_blank"
			class="hover:underline">Pricing</a
		>
		<a href="https://github.com/humanbeeng/checkpost">
			<GithubLogo />
		</a>
	</Header>
	<img
		src={background}
		alt=""
		class="hidden md:block absolute inset-0 -z-10 h-full w-full stroke-gray-200 [mask-image:radial-gradient(100%_100%_at_top_right,white,transparent)]"
	/>

	<div class=" grid grid-cols-1 mx-6 md:grid-cols-2 md:max-w-7xl lg:mx-auto">
		<div class="w-auto flex my-32 md:my-64 flex-col md:justify-center">
			<h1
				class="text-4xl antialiased font-bold leading-tight select-none tracking-tighter bg-gradient-to-r text-left lg:text-6xl lg:text-left lg:leading-[1.1]"
			>
				Monitor Incoming Webhooks in Realtime
			</h1>
			<section>
				<div class="flex flex-col pt-4">
					<div class="flex place-items-center gap-2">
						<div
							class="border place-items-center inline-flex justify-between border-black rounded-sm h-10"
						>
							<input
								type="text"
								class="ml-2 border-none h-9 outline-none w-28 on:click:outline-none place-self-center bg-inherit"
								placeholder="boringcorp"
								size="10"
								spellcheck="false"
								maxlength="10"
								minlength="4"
								name="subdomain"
								autocomplete="off"
								value={subdomain}
								on:keydown={handleInput}
								pattern="[A-Za-z0-9]+"
							/>
							<p
								class="font-normal text-md tracking-tight pr-3 place-self-center text-left justify-self-start"
							>
								.checkpost.io
							</p>
						</div>
						<Button variant="default" class="h-10" disabled={state !== 'success'} on:click={signin}>
							{#if state === 'loading'}
								<Reload class="h-4  animate-spin" />
							{:else}
								<ChevronRight />
							{/if}
						</Button>
					</div>
					<div>
						{#if state === 'empty' || state === 'loading'}
							<h3
								class="font-light text-md opacity-60 select-none leading-10 md:leading-10 text-left"
							>
								Claim Your Unique Subdomain - For You or Your Team, Absolutely Free
							</h3>
						{/if}
						{#if state === 'success' && existsRes}
							{#if existsRes && !existsRes.exists}
								<p class="text-green-900 py-2">{existsRes?.message}</p>
							{:else if existsRes}
								<p class="text-red-500 py-2">{existsRes.message}</p>
							{/if}
						{/if}
						{#if state === 'error'}
							<p class="text-red-500 py-2">{error}</p>
						{/if}
					</div>
				</div>
			</section>
		</div>
	</div>
	<Footer />
</div>
