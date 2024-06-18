<script lang="ts" type="module">
	import { page } from '$app/stores';
	import { PUBLIC_WEBSOCKET_URL } from '$env/static/public';
	// import { PUBLIC_WEBSOCKET_URL } from '$env/static/public';
	import logo from '$lib/assets/logo-black.svg';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import HistoryItem from '@/components/HistoryItem.svelte';
	import MethodBadge from '@/components/MethodBadge.svelte';
	import ProfileHeader from '@/components/ProfileHeader.svelte';
	import RequestContainer from '@/components/RequestContainer.svelte';
	import RequestDetails from '@/components/RequestDetails.svelte';
	import StatusCodeBadge from '@/components/StatusCodeBadge.svelte';
	import { Button } from '@/components/ui/button';
	import { endpointHistory } from '@/store.js';
	import type { Request, WebsocketPayload } from '@/types.js';
	import clsx from 'clsx';
	import ReconnectingWebSocket from 'reconnecting-websocket';
	import ResilientWebSocket, { WebSocketEvent } from 'resilient-websocket';

	import { onMount } from 'svelte';

	import { Exit } from 'svelte-radix';

	const endpoint = $page.params.endpoint;
	export let data;

	let selectedRequest: Request | undefined;

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

	let socket: ReconnectingWebSocket;
	const connectSocket = () => {
		const wsUrl = `${PUBLIC_WEBSOCKET_URL}/endpoint/inspect/${endpoint}?token=${data.token}`;

		socket = new ReconnectingWebSocket(wsUrl);

		// Connection opened
		socket.addEventListener('open', (event) => {
			console.log('Websocket connection established');
		});

		// Listen for messages
		socket.addEventListener('message', (event) => {
			//ping
			if (event.data === '') {
				return;
			}

			const req: WebsocketPayload = JSON.parse(event.data);
			if (req.code == 200) {
				$endpointHistory.requests = [req.hook_request, ...($endpointHistory.requests ?? [])];
			}
		});

		socket.addEventListener('close', () => {
			console.log('Closing');
			socket.close();
		});
	};

	onMount(() => {
		connectSocket();
		// Send ping

		setInterval(() => {
			if (socket && socket.readyState === socket.OPEN) {
				socket.send('');
			}
		}, 5000);
	});
</script>

<body class="bg-gray-50 flex overflow-hidden h-screen">
	<!-- Sidebar -->
	<div class="min-w-64 max-w-64 border-r border-gray-300 bg-gray-200 flex flex-col justify-between">
		<!-- Branding -->
		<div class="border-b border-gray-300 px-5 py-4">
			<span class="flex">
				<img src={logo} alt="Checkpost logo" />
				<p class=" tracking-normal font-medium text-md">Checkpost</p>
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

		<!-- User button -->
		<DropdownMenu.Root>
			<DropdownMenu.Trigger asChild let:builder>
				<Button
					variant="ghost"
					builders={[builder]}
					class="rounded-md  flex gap-1 justify-start  bg-gray-100  border border-gray-300 m-4 px-4 py-6 shadow-sm hover:bg-gray-50"
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
	<div class=" flex-1 overflow-y-auto w-screen">
		<!-- Header -->
		{#if selectedRequest}
			<div class="flex justify-between mt-1 py-4 px-10 border-b border-gray-300 gap-4">
				<span class="flex gap-2 w-3/4">
					<MethodBadge method={selectedRequest.method} />
					{#if selectedRequest.path === '/'}
						<p class="">
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
		<div class="my-4 mx-10 flex flex-col gap-2 overflow-y-auto">
			{#if selectedRequest}
				<!-- <RequestStats successCount={100} /> -->
				<RequestDetails request={selectedRequest} />
			{:else}
				<div class="flex flex-col justify-start w-full my-32">
					<p class="text-3xl font-bold tracking-tight my-2">It's empty in here</p>
					<span class="text-lg font-normal text-gray-800"
						>Try calling this endpoint
						<code class="bg-gray-50 py-1 border px-4 rounded-md">
							<a href="http://{endpoint}.checkpost.local:3000" class="underline" target="_blank">
								https://{endpoint}.checkpost.io/
							</a>
						</code>
					</span>
				</div>
			{/if}
		</div>
	</div>
</body>
