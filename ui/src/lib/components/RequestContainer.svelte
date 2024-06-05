<script lang="ts">
	import type { Request } from '@/types';

	export let request: Request;
</script>

<div class="flex flex-col outline gap-2 m-4">
	<p>UUID: {request.uuid}</p>
	<p>Path: {request.path}</p>
	<p>Headers: {request.headers}</p>
	<p>Method: {request.method}</p>
	<p>Response code: {request.response_code}</p>
	<p>IP: {request.source_ip}</p>
	<pre id="json">Content: {request.content}</pre>
	<p>Content Size: {request.content_size}</p>
</div>
