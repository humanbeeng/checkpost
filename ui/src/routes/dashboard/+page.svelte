<script>
	import Header from '@/components/Header.svelte';

	export let data;
</script>

<body class="h-screen bg-gray-50">
	<Header user={data.user ?? null} />
	<h1>Dashboard</h1>
</body>
