import { error } from '@sveltejs/kit';
import type { PageServerLoad } from '../../$types';

export const load: PageServerLoad = async ({ fetch, cookies }) => {

	const res = await fetch('http://localhost:3000/admin/dashboard', {
		method: 'GET'
	})

	if (res.ok) {
		const data = await res.json()
		console.log(data)
		return data
	} else {
		return error(400)
	}
};
