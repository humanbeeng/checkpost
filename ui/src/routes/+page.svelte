<script lang="ts">
	import background from '$lib/assets/background.svg';
</script>

<nav>
	<a href="/">Checkpost</a>
	<a href="/docs">Docs</a>
</nav>

<div style="background-image: url({background}); " class="bg-cover w-full h-screen">
	<div class=" grid grid-cols-1 mx-6 h-screen md:grid-cols-2 md:mx-10">
		<div class="w-auto flex my-10 flex-col md:justify-center">
			<h1 class="text-5xl font-bold text-left">Inspect HTTP requests in realtime</h1>
		</div>
		<div class="hidden md:w-auto md:block">hello</div>
	</div>
</div>
