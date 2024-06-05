<script lang="ts">
	export let requestCount = 0;
	export let successCount = 0;
	export let failureCount = 0;
</script>

<div class="w-full flex gap-2 justify-between border bg-gray-100 px-8 py-1 rounded-md shadow-sm">
	<p>Total count: {requestCount}</p>
	<p>Success count: {successCount}</p>
	<p>Failure count: {successCount}</p>
</div>
