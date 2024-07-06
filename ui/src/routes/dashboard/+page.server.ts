import { PUBLIC_BASE_URL } from '$env/static/public';
import type { User } from '@/types';
import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from '../$types';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
	if (!cookies.get('token')) {
		return redirect(301, '/auth/logout');
	}

	const res = await fetch(`${PUBLIC_BASE_URL}/user`).catch((err) => {
		console.log('Unable to fetch user details', err);
		throw error(500);
	});

	if (!res.ok) {
		// TODO: Better error handling
		console.log('err', await res.text());

		if (res.status == 401) {
			return redirect(301, '/auth/logout');
		}
		return {
			err: { message: 'Unable to fetch user details' }
		};
	}

	const user = (await res.json().catch((err) => {
		console.error('Unable to parse user response', err);
	})) as User;

	return { user };
};
