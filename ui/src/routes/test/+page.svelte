<script lang="ts">
	import type { PageData } from './$types';

	export let data: PageData;
</script>

<body class="w-full h-screen overflow-hidden flex">
	<div class="flex flex-col bg-gray-200 border-r">
		<div class="h-16">
			<a href="/">Checkpost</a>
		</div>
	</div>

	<div class="flex-1 flex-col overflow-y-auto gap-1">
		<h1>Main content</h1>

		<div class="h-screen w-screen bg-green-100">Content</div>
		<div class="h-screen w-screen bg-green-100">Content</div>
		<div class="h-screen w-screen bg-green-100">Content</div>
		<div class="h-screen w-screen bg-green-100">Content</div>
	</div>
</body>
