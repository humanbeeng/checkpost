<script lang="ts">
	import background from '$lib/assets/background.svg';
	import { Button } from '$lib/components/ui/button/index.js';
</script>

<header class="sticky border-b border-black/10 w-full bg-[#D9D9D9]">
	<div class="container flex h-12 max-w-screen-2xl items-center justify-between">
		<a href="/" class="font-medium">Checkpost</a>
		<nav class="flex items-center gap-6 text-sm">
			<a href="/pricing" class="hover:underline">Pricing</a>
			<a href="/about" class="hover:underline">About</a>
			<a href="/login"><Button class="h-8 border-black" variant="default">Sign in</Button></a>
		</nav>
	</div>
</header>
<div style="background-image: url({background});" class="bg-cover w-full min-h-screen">
	<div class=" grid grid-cols-1 mx-10 min-h-screen md:grid-cols-2 md:max-w-7xl lg:mx-auto">
		<div class="w-auto flex my-32 flex-col md:justify-center">
			<h1
				class="text-4xl font-bold leading-tight select-none tracking-tighter bg-gradient-to-r text-left lg:text-6xl lg:text-left lg:leading-[1.1]"
			>
				Inspect incoming HTTP requests in realtime
			</h1>

			<section>
				<h3
					class="font-light text-xl opacity-60 pt-2 select-none leading-10 md:py-4 md:leading-10 text-left"
				>
					Within seconds, you'll get
					<ul class="list-disc mx-5 font-normal text-xl">
						<li>Access to a free endpoint</li>
						<li>View incoming requests in realtime</li>
					</ul>
				</h3>
			</section>

			<Button class="w-full h-11 bg-foreground my-4 lg:w-1/2">Coming soon.</Button>
		</div>
		<div class="hidden md:w-auto md:block"></div>
	</div>
</div>
