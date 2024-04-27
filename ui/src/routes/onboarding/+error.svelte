<script lang="ts">
	import { page } from '$app/stores';
</script>

<body>
	<h1>{$page.error?.message}</h1>
</body>
