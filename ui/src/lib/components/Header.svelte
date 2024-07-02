<script lang="ts">
	import logo from '$lib/assets/logo-black.svg';
	import type { User } from '@/types';
	import ProfileHeader from './ProfileHeader.svelte';
	export let user: User | null;
</script>

<header
	class="sticky top-0 border-b border-black/10 w-full bg-gray-100/10 backdrop-blur-3xl z-10 shadow-sm"
>
	<div class="container flex h-12 max-w-screen-2xl items-center justify-between">
		<a href="/" class="font-medium flex gap-1"
			><span>
				<img src={logo} alt="Checkpost logo" />
			</span>
			<p class=" tracking-tight">Checkpost</p></a
		>

		<nav class="flex items-center gap-6 text-sm">
			<slot />
			{#if user}
				<ProfileHeader {user} />
			{/if}
		</nav>
	</div>
</header>
