<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import Header from '@/components/Header.svelte';
	import Reload from 'svelte-radix/Reload.svelte';

	import { enhance } from '$app/forms';
	import { user } from '$lib/store';
	import * as Avatar from '@/components/ui/avatar';
	import { Button } from '@/components/ui/button';
	import { debounce } from '@/debounce';
	import { Link1, LinkBreak2 } from 'svelte-radix';
	import type { EndpointExistsResponse, State } from '../types';
	import type { User } from './types';

	export let data: User;
	if (data) {
		user.set(data);
	}
	// export let form;

	let error: string | null;
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

		const res = await fetch(`http://api.checkpost.local:3000/url/exists/${subdomain}`);

		switch (res.status) {
			case 200: {
				endpointExistsResponse = (await res.json()) as EndpointExistsResponse;
				error = null;
				state = 'success';
				break;
			}
			case 400: {
				error = await res.text();
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
		}
		debouncedHandleChange(e);
	}
</script>

<body class="h-screen flex flex-col">
	<Header>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger asChild let:builder>
				<Button variant="ghost" builders={[builder]} class="relative  rounded-md">
					<Avatar.Root class="h-8 w-8">
						<Avatar.Image src="https://placehold.co/32x32.png" alt={data.name} />
						<Avatar.Fallback>NR</Avatar.Fallback>
					</Avatar.Root>
					<p class="px-2" autocapitalize="on">{$user?.name}</p>
				</Button>
			</DropdownMenu.Trigger>
			<DropdownMenu.Content class="w-56" align="end">
				<DropdownMenu.Item>Log out</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Header>

	<main class="w-full items-center flex flex-col flex-grow justify-center bg-gray-100/10">
		<div
			class="border py-3 px-4 rounded-xl mb-52 flex flex-col items-center justify-center gap-6 bg-white/90 shadow-md hover:isolate"
		>
			<div class="w-full">
				<h3 class="self-start font-medium text-lg px-1">Create your endpoint</h3>
				<hr class="mt-1" />
			</div>
			<div class="flex gap-2 p-1">
				<form method="POST" use:enhance>
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
									<p class="text-red-900 text-xs py-2">{error}</p>
								{/if}
							</span>
						</div>
						<!-- <label for="subdomain">Your URL</label> -->
						<Button class="w-full" disabled={state !== 'success'} variant="default" type="submit">
							{#if state === 'loading'}
								<Reload class="mr-2 h-4 w-4 animate-spin" />
								Checking
							{:else}
								Continue
							{/if}
						</Button>
					</div>
				</form>
			</div>
		</div>
	</main>
</body>
