<script lang="ts" type="module">
	import { page } from '$app/stores';
	import { PUBLIC_WEBSOCKET_URL } from '$env/static/public';
	import logo from '$lib/assets/logo-black.svg';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import HistoryItem from '@/components/HistoryItem.svelte';
	import MethodBadge from '@/components/MethodBadge.svelte';
	import OnlineStatusIndicator from '@/components/OnlineStatusIndicator.svelte';
	import RequestDetails from '@/components/RequestDetails.svelte';
	import StatusCodeBadge from '@/components/StatusCodeBadge.svelte';
	import { Button } from '@/components/ui/button';
	import { endpointHistory } from '@/store.js';
	import type { Request, WSMessage } from '@/types.js';
	import clsx from 'clsx';

	import { onMount } from 'svelte';

	import { Exit, Reload } from 'svelte-radix';

	export let data;

	const user = data.user;
	const endpoint = $page.params.endpoint;
	let selectedRequest: Request | undefined;
	let websocketOnline: 'connecting' | 'offline' | 'online' = 'connecting';

	$endpointHistory = data.endpointHistory;
	if ($endpointHistory == null) {
		console.log('Endpoint History is null');
	} else {
		if ($endpointHistory.requests?.length) {
			selectedRequest = $endpointHistory.requests.at(0);
		}
	}

	const selectRequest = (requestuuid: string) => {
		selectedRequest = $endpointHistory?.requests?.find((r) => r.uuid == requestuuid);
	};

	const connectSocket = () => {
		const wsUrl = `${PUBLIC_WEBSOCKET_URL}/endpoint/inspect/${endpoint}?token=${data.token}`;
		const socket = new WebSocket(wsUrl);

		// Connection opened
		socket.addEventListener('open', () => {
			console.log('Websocket connection established');
			websocketOnline = 'online';
		});

		// Listen for messages
		socket.addEventListener('message', (event) => {
			const message = JSON.parse(event.data) as WSMessage;
			if (message.code == 200) {
				$endpointHistory.requests = [message.payload, ...($endpointHistory.requests ?? [])];
			} else {
				// TODO: Handle error
			}
		});

		socket.addEventListener('close', () => {
			console.log('Websocket connection closed');
			websocketOnline = 'offline';
		});

		socket.addEventListener('error', () => {
			console.log('Websocket connection error');
			websocketOnline = 'offline';
		});
	};

	onMount(() => {
		connectSocket();
	});
</script>

<body class="bg-gray-50 flex overflow-hidden h-screen">
	<!-- Sidebar -->
	<div class="min-w-64 max-w-64 border-r border-gray-300 bg-gray-200 flex flex-col justify-between">
		<!-- Branding -->
		<div class="border-b border-gray-300 px-5 py-4">
			<span class="flex justify-between">
				<span class="flex gap-1 place-items-center">
					<img src={logo} alt="Checkpost logo" />
					<p class=" tracking-normal font-medium text-md">Checkpost</p>
				</span>
				{#if websocketOnline === 'offline'}
					<Button
						class="h-7 w-auto p-2 bg-red-100 text-red-600 hover:bg-red-200"
						on:click={connectSocket}
					>
						<span class="flex gap-1 place-items-center">
							<Reload class="size-3" />Reconnect
						</span>
					</Button>
				{:else}
					<OnlineStatusIndicator state={websocketOnline} />
				{/if}
			</span>
		</div>

		<!-- History -->
		<div
			class="px-5 flex flex-col justify-self-start grow overflow-y-auto border-b border-b-gray-300"
		>
			<p class="font-medium text-md my-4 text-gray-600">Request history</p>
			{#if $endpointHistory && $endpointHistory.requests}
				{#each $endpointHistory.requests as request}
					<button
						on:click={() => selectRequest(request.uuid)}
						class={clsx(
							'rounded-md',
							'px-2',
							selectedRequest?.uuid == request.uuid && 'bg-gray-300 shadow-sm'
						)}
					>
						<HistoryItem {request} />
					</button>
				{/each}
			{:else}
				<p class="font-light text-gray-400">No requests</p>
			{/if}
		</div>

		{#if user.plan === 'free'}
			<div class="px-2 mx-4 py-2 bg-orange-200/90 mt-2 rounded-md">
				<p class="font-extralight text-gray-600 text-xs">
					⚠️ You're on {user.plan} tier. <br />Incoming requests history will be expired in 6 hours.
				</p>
			</div>
		{/if}

		<!-- User button -->
		<DropdownMenu.Root>
			<DropdownMenu.Trigger asChild let:builder>
				<Button
					variant="ghost"
					builders={[builder]}
					class="rounded-md  flex gap-1 justify-start  bg-gray-100  border border-gray-300 mx-4 mt-2 mb-3 px-4 py-6 shadow-sm hover:bg-gray-50"
				>
					<img src={data.user.avatar_url} alt={data.user.name} class="h-8 rounded-md" />
					<p class="px-2 text-md" autocapitalize="on">{data.user.name}</p>
				</Button>
			</DropdownMenu.Trigger>
			<DropdownMenu.Content class="w-72" sameWidth={true}>
				<a href="/auth/logout" data-sveltekit-reload>
					<DropdownMenu.Item>
						<span class="flex gap-1 justify-start align-middle">
							<Exit class="h-4" /> Log out
						</span>
					</DropdownMenu.Item>
				</a>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</div>

	<!-- Main section -->
	<div class="flex-1 overflow-y-auto w-screen">
		<!-- Header -->
		{#if selectedRequest}
			<div
				class="flex justify-between py-4 px-10 border-b border-gray-300 gap-4 sticky top-0 backdrop-blur-md"
			>
				<span class="flex gap-2 w-3/4">
					<MethodBadge method={selectedRequest.method} />
					{#if selectedRequest.path === '/'}
						<p>
							{selectedRequest.path}
						</p>
					{:else}
						<p class="truncate">
							{selectedRequest.path}
						</p>
					{/if}
				</span>
				<StatusCodeBadge response_code={selectedRequest.response_code} />
			</div>
		{/if}
		<!-- Request details -->
		<div class="my-4 mx-10 flex flex-col gap-2 h-full">
			{#if selectedRequest}
				<RequestDetails request={selectedRequest} />
			{:else}
				<div class="flex flex-col justify-start w-full my-32">
					<p class="text-3xl font-bold tracking-tight my-2">It's empty in here</p>
					<span class="text-lg font-normal text-gray-800"
						>Try calling this endpoint
						<code class="bg-gray-50 py-1 border px-4 rounded-md">
							<a href="http://{endpoint}.checkpost.io" class="underline" target="_blank">
								https://{endpoint}.checkpost.io/
							</a>
						</code>
					</span>
				</div>
			{/if}
		</div>
	</div>
</body>
