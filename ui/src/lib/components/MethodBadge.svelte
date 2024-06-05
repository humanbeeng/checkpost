<script lang="ts">
	import type { HttpMethod } from '@/types';
	import clsx from 'clsx';
	export let method: HttpMethod;
</script>

<div
	class={clsx(
		method.toLowerCase() == 'get' && 'bg-indigo-600',
		method.toLowerCase() == 'post' && 'bg-green-600',
		method.toLowerCase() == 'put' && 'bg-orange-600',
		method.toLowerCase() == 'delete' && 'bg-red-600',

		method.toLowerCase() == 'patch' && 'bg-slate-600',
		method.toLowerCase() == 'head' && 'bg-pink-600',
		method.toLowerCase() == 'options' && 'bg-teal-600',
		'font-bold',
		'w-fit',
		'inline-flex',
		'items-center',
		'px-2',
		'py-1',
		'rounded-md',
		'text-secondary',
		'text-xs',
		'max-h-6'
	)}
>
	<p>{method.toUpperCase()}</p>
</div>
