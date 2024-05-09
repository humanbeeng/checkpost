<script lang="ts">
	import { goto } from '$app/navigation';
	import background from '$lib/assets/background.svg';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import Header from '@/components/Header.svelte';

	import { PUBLIC_SERVER_URL } from '$env/static/public';
	import { debounce } from '@/debounce';
	import { Check, ChevronRight, Cross2, Reload } from 'svelte-radix';
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

	<section class="h-screen border-t bg-gray-50 lg:px-24">
		<h1 class="text-3xl font-bold tracking-tight text-center py-10 lg:py-12">Plans</h1>
		<div class="w-full px-6 lg:px-32 grid grid-cols-1 gap-3 lg:grid-cols-3">
			<!-- Free -->
			<Card.Root class="mx-2 ">
				<Card.Header>
					<Card.Title tag="h1" class="text-xl">Free</Card.Title>
					<Card.Description>Signed in users</Card.Description>
				</Card.Header>
				<Card.Content>
					<ul class="flex flex-col gap-2 list-none mx-5">
						<li><span class="flex gap-2"><Check />1 custom endpoint</span></li>
						<li><span class="flex gap-2"><Check />50 incoming requests/day</span></li>
						<li><span class="flex gap-2"><Check />4 hour request history retention</span></li>
						<li><span class="flex gap-2"><Check />Port forwarding (Coming soon)</span></li>
						<li><span class="flex gap-2"><Check />File attachments upto 10MB</span></li>
						<li><span class="flex gap-2"><Check />Request payload limit 10KB</span></li>
						<li><span class="flex gap-2"><Check />Rate limit: 3 rps</span></li>
						<li><span class="flex gap-2"><Check />Discord access - Community</span></li>
						<li>
							<span class="flex gap-2"
								><Cross2 />
								<p>Proxy middleware</p></span
							>
						</li>
						<li>
							<span class="flex gap-2"
								><Cross2 />
								<p>Password protected endpoint</p></span
							>
						</li>
					</ul>
				</Card.Content>
			</Card.Root>

			<!-- Basic -->
			<Card.Root class="outline-primary outline">
				<Card.Header>
					<Card.Title tag="h1" class="text-xl">Basic</Card.Title>
					<Card.Description
						><span class="font-medium text-black underline">$49/year</span>
						or <span class="font-medium text-xs">$5/mo</span></Card.Description
					>
				</Card.Header>
				<Card.Content>
					<ul class="mx-5 list-none flex flex-col gap-2">
						<li><span class="flex gap-2"><Check />1 custom endpoint</span></li>
						<li><span class="flex gap-2"><Check />Unlimited incoming reqs/day</span></li>
						<li><span class="flex gap-2"><Check />30 day request history retention</span></li>
						<li><span class="flex gap-2"><Check />Port forwarding (Coming soon)</span></li>
						<li>
							<span class="flex gap-2"><Check />File attachments upto 20MB (Coming soon)</span>
						</li>
						<li><span class="flex gap-2"><Check />Request payload limit 510Kb</span></li>
						<li><span class="flex gap-2"><Check />Rate limit: 10 rps</span></li>
						<li><span class="flex gap-2"><Check />Discord access - Private</span></li>
						<li><span class="flex gap-2"><Check />1 proxy middleware</span></li>
						<li><span class="flex gap-2"><Check />Password protected endpoint</span></li>
					</ul>
				</Card.Content>
				<Card.Footer>
					<Button class="w-full bg-black">Coming soon</Button>
				</Card.Footer>
			</Card.Root>

			<!--Pro-->
			<Card.Root class="mx-2 ">
				<Card.Header>
					<Card.Title tag="h1" class="text-xl">Team</Card.Title>

					<Card.Description
						><span class="font-medium text-md text-black">$49/year</span>
						or <span class="font-medium text-xs">$5/mo</span></Card.Description
					>
				</Card.Header>
				<Card.Content>
					<ul class="list-disc mx-5">
						<li>5 custom endpoints.</li>
						<li>Request history is stored upto 180 days.</li>
						<li>Auto cleanup of request history (Opt in)</li>
						<li>Port forwarding</li>
						<li>Request payload limit 510Kb.</li>
						<li>Sharable request details.</li>
						<li>Upto 5 collaborators.</li>
						<li>Password protected endpoints.</li>
						<li>Unlimited middlewares. With freeze and edit before forward request.</li>
						<li>Rate limit: 15 rps.</li>
						<li>Discord access - Private Channel.</li>
						<li>Priority support.</li>
						<li>File attachments</li>
						<li>Font switcher for code view.</li>
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
