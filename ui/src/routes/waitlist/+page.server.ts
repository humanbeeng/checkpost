import { PUBLIC_BASE_URL } from '$env/static/public';
import type { User } from '@/types.js';
import { error, redirect } from '@sveltejs/kit';
import type { UserEndpointsResponse } from '../onboarding/types.js';

export const load = async ({ fetch, cookies }) => {
	const token = cookies.get('token');
	if (!token) {
		return redirect(301, '/');
	}

	const fetchUser = async () => {
		const res = await fetch(`${PUBLIC_BASE_URL}/user`).catch((err) => {
			console.error('Unable to fetch user details', err);
			return error(500);
		});
		if (!res.ok) {
			// TODO: Better error handling
			return error(res.status, await res.text());
		}

		const user = (await res.json().catch((err) => {
			console.error('Unable to parse user response', err);
		})) as User;

		return user;
	};

	const fetchEndpoints = async () => {
		const res = await fetch(`${PUBLIC_BASE_URL}/endpoint`).catch((err) => {
			console.log('Unable to fetch user details', err);
			return error(500);
		});

		if (!res.ok) {
			// TODO: Better error handling
			return error(res.status, await res.text());
		}

		const endpoints = (await res.json().catch((err) => {
			console.error('Unable to parse user response', err);
		})) as UserEndpointsResponse;

		if (!endpoints.endpoints || !endpoints.endpoints.length) {
			return redirect(301, '/onboarding');
		}

		return endpoints;
	};

	return { user: await fetchUser(), endpoints: await fetchEndpoints(), err: null };
};
