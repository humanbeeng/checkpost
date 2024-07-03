<script lang="ts">
	import Header from '@/components/Header.svelte';
	import { Button } from '@/components/ui/button';
</script>

<div class="h-screen bg-gray-50">
	<Header user={null} />

	<div class="flex flex-col w-full items-center justify-center mt-72">
		<!-- Box  -->
		<div
			class="rounded-md w-1/2 mx-20 border flex flex-col p-6 gap-3 justify-between xl:w-1/5 shadow-lg backdrop-blur-lg bg-gray-50"
		>
			<div>
				<h1 class="text-xl font-medium">Login</h1>
				<h2 class="font-normal text-md tracking-tight text-gray-600">to continue to Checkpost</h2>
			</div>
			<hr />
			<div class="flex flex-col gap-3 pt-2">
				<a href="/auth/github" class="w-full">
					<Button variant="default" class="w-full">Continue with Github</Button>
				</a>
				<a href="/auth/google" class="w-full">
					<Button variant="default" class="w-full">Continue with Google</Button>
				</a>
			</div>
		</div>
	</div>
</div>
