<script lang="ts">
	import { goto } from '$app/navigation';
	import background from '$lib/assets/background.svg';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import Header from '@/components/Header.svelte';

	import { PUBLIC_SERVER_URL } from '$env/static/public';
	import { debounce } from '@/debounce';
	import { ChevronRight, Reload } from 'svelte-radix';
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

		const res = await fetch(`${PUBLIC_SERVER_URL}/url/exists/${subdomain}`);

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

<Header user={null}>
	<a href="/pricing" class="hover:underline">Pricing</a>
</Header>

<div class=" w-full min-h-full">
	<img
		src={background}
		alt=""
		class="absolute inset-0 -z-10 h-full w-full stroke-gray-200 [mask-image:radial-gradient(100%_100%_at_top_right,white,transparent)]"
	/>

	<div class=" grid grid-cols-1 mx-6 md:min-h-screen md:grid-cols-2 md:max-w-7xl lg:mx-auto">
		<div class="w-auto flex my-32 flex-col md:justify-center">
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
					<div class="">
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
		<!--Dashboard image-->
		<div class="hidden md:w-auto md:block"></div>
	</div>

	<section class="h-screen border-t bg-gray-50">
		<h1 class="text-3xl font-bold tracking-tight text-center py-10 md:py-12">Plans</h1>
		<div class="w-full px-6 lg:px-32 grid grid-cols-1 gap-3 md:grid-cols-4">
			<!-- Guest -->
			<Card.Root>
				<Card.Header>
					<Card.Title tag="h1" class="text-xl">Guest</Card.Title>
					<Card.Description>Non signed in users</Card.Description>
				</Card.Header>
				<Card.Content>
					<ul class="list-disc mx-5">
						<li>
							One random endpoint. <br />Example: <a href="/"><u>https://xyzxa.checkpost.io</u></a>
						</li>
						<li>Expires in 1 day.</li>
						<li>Request history retained for 1 day.</li>
						<li>Request payload limit 32Kb.</li>
						<li>Rate limit 2 rps.</li>
					</ul>
				</Card.Content>
			</Card.Root>

			<!-- Free -->
			<Card.Root>
				<Card.Header>
					<Card.Title tag="h1" class="text-xl">Guest</Card.Title>
					<Card.Description>Non signed in users</Card.Description>
				</Card.Header>
				<Card.Content>
					<ul class="list-disc mx-5">
						<li>
							One random endpoint. <br />Example: <a href="/"><u>https://xyzxa.checkpost.io</u></a>
						</li>
						<li>Expires in 1 day.</li>
						<li>Request history retained for 1 day.</li>
						<li>Request payload limit 32Kb.</li>
						<li>Rate limit 2 rps.</li>
					</ul>
				</Card.Content>
			</Card.Root>

			<!-- Basic -->
			<Card.Root>
				<Card.Header>
					<Card.Title tag="h1" class="text-xl">Basic</Card.Title>
					<Card.Description>$4/month</Card.Description>
				</Card.Header>
				<Card.Content>
					<ul class="list-disc mx-5">
						<li>
							One random endpoint. <br />Example: <a href="/"><u>https://xyzxa.checkpost.io</u></a>
						</li>
						<li>Expires in 2 days.</li>
						<li>Request history retained for 2 days.</li>
						<li>Request payload limit 32Kb.</li>
						<li>Rate limit: 3 rps.</li>
						<li>Discord Access - Community channel.</li>
					</ul>
				</Card.Content>
				<!-- <Card.Footer> -->
				<!-- 	<Button class="w-full">Generate endpoint</Button> -->
				<!-- </Card.Footer> -->
			</Card.Root>

			<!--Pro-->
			<Card.Root>
				<Card.Header>
					<Card.Title tag="h1" class="underline text-xl">Pro</Card.Title>
					<Card.Description>$5/month | $55/year</Card.Description>
				</Card.Header>
				<Card.Content>
					<ul class="list-disc mx-5">
						<li>
							10 custom endpoints. <br /><a href="/"
								><u>Example: https://yourname.checkpost.io</u></a
							>
						</li>
						<li>Request history is stored upto 180 days.</li>
						<li>10 Web CRONs.</li>
						<li>Request payload limit 510Kb.</li>
						<li>Password protected endpoints.</li>
						<li>File attachments support.</li>
						<li>Port forwarding for local servers.</li>
						<li>1 proxy middleware.</li>
						<li>Rate limit: 30 rps.</li>
						<li>Discord Access - Private Channel</li>
					</ul>
				</Card.Content>
				<Card.Footer>
					<Button class="w-full bg-black underline">Upgrade - Coming soon</Button>
				</Card.Footer>
			</Card.Root>

			<!--Enterprise-->
			<!-- <Card.Root> -->
			<!-- 	<Card.Header> -->
			<!-- 		<Card.Title tag="h1" class="text-xl">Enterprise (Coming soon)</Card.Title> -->
			<!-- 		<Card.Description>Enterprise (Coming soon)</Card.Description> -->
			<!-- 	</Card.Header> -->
			<!-- <Card.Content> -->
			<!-- 	<ul class="list-disc mx-5"> -->
			<!-- 		<li>Inspect incoming HTTP requests.</li> -->
			<!-- 		<li>1 endpoint</li> -->
			<!-- 		<li>Request history retained for 2 days</li> -->
			<!-- 		<li>Expires in 4h.</li> -->
			<!-- 		<li>Request payload limit 32Kb.</li> -->
			<!-- 		<li>Rate limit: 10 rps.</li> -->
			<!-- 		<li>Discord Access</li> -->
			<!-- 	</ul> -->
			<!-- </Card.Content> -->
			<!-- <Card.Footer> -->
			<!-- 	<Button class="w-full">Generate endpoint</Button> -->
			<!-- </Card.Footer> -->
			<!-- </Card.Root> -->
		</div>
		<footer class="py-4 md:px-8 md:py-0 bottom-0">
			<div class="container flex flex-col items-center justify-between gap-4 md:h-24 md:flex-row">
				<div class="flex flex-col items-center gap-4 px-8 md:flex-row md:gap-2 md:px-0">
					<p class="text-center text-sm leading-loose text-muted-foreground md:text-left">
						Built by <a
							href="https://twitter.com/nithinrajx"
							target="_blank"
							rel="noreferrer"
							class="font-medium underline underline-offset-4"
							data-svelte-h="svelte-18aekvv">Nithin Raj</a
						>.
					</p>
				</div>
			</div>
		</footer>
	</section>
</div>
