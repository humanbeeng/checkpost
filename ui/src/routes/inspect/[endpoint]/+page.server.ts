import { PUBLIC_BASE_URL } from '$env/static/public';
import type { EndpointHistory, User } from '@/types';
import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const csr = true;

export const load: PageServerLoad = async ({ fetch, params, cookies }) => {
	const endpoint = params.endpoint;
	const token = cookies.get('token');

	const fetchEndpointHistory = async () => {
		const res = await fetch(`${PUBLIC_BASE_URL}/endpoint/history/${endpoint}`).catch((err) => {
			console.error('Unable to fetch endpoint request history', err);
			error(500);
		});

		if (!res.ok) {
			if (res.status == 401) {
				throw redirect(301, '/auth/logout');
			}
			error(res.status, { message: await res.text() });
		}

		const endpointHistory = (await res.json().catch((err) => {
			console.error('Unable to parse endpoint history', err);
		})) as EndpointHistory;

		return endpointHistory;
	};

	const fetchUser = async () => {
		console.log('Fetching user details');
		const res = await fetch(`${PUBLIC_BASE_URL}/user`).catch((err) => {
			console.error('Unable to fetch user details', err);
			error(500);
		});

		if (!res.ok) {
			if (res.status == 401) {
				throw redirect(301, '/auth/logout');
			} else if (res.status == 404) {
				throw redirect(301, '/onboarding');
			} else if (res.status === 403) {
				error(res.status, { message: await res.text() });
			}
		}

		const user = (await res.json().catch((err) => {
			console.error('Unable to parse user response', err);
			error(500, { message: 'Something went wrong' });
		})) as User;

		return user;
	};

	const user = await fetchUser();
	const endpointHistory = await fetchEndpointHistory();

	return {
		user,
		endpointHistory,
		token,
	};
};
