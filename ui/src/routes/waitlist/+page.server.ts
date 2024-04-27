import { PUBLIC_SERVER_URL } from '$env/static/public';
import type { User } from '@/types.js';
import { error, redirect } from '@sveltejs/kit';

export const load = async ({ fetch, cookies }) => {
	const token = cookies.get('token');
	if (!token) {
		return redirect(301, '/');
	}

	const res = await fetch(`${PUBLIC_SERVER_URL}/user`).catch((err) => {
		console.log('Unable to fetch user details', err);
		return error(500);
	});

	if (!res.ok) {
		// TODO: Better error handling
		return {
			err: { message: 'Unable to fetch user details' }
		};
	}

	const user = (await res.json().catch((err) => {
		console.error('Unable to parse user response', err);
	})) as User;

	console.log('User fetched from server', user);
	return { user };
};
