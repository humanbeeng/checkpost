<script lang="ts">
	import { user } from '$lib/store';
</script>

{#if $user}
	hello
{:else}
	<p>Not signed in</p>
{/if}
