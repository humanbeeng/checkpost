import { PUBLIC_SERVER_URL } from '$env/static/public';
import type { User } from '@/types.js';
import { error, redirect } from '@sveltejs/kit';
import type { UserEndpointsResponse } from '../onboarding/types.js';

export const load = async ({ fetch, cookies }) => {
	const token = cookies.get('token');
	if (!token) {
		return redirect(301, '/');
	}

	const fetchUser = async () => {
		const res = await fetch(`${PUBLIC_SERVER_URL}/user`).catch((err) => {
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

	const fetchUrls = async () => {
		const res = await fetch(`${PUBLIC_SERVER_URL}/url`).catch((err) => {
			console.log('Unable to fetch user details', err);
			return error(500);
		});

		if (!res.ok) {
			// TODO: Better error handling
			return error(res.status, await res.text());
		}

		const urls = (await res.json().catch((err) => {
			console.error('Unable to parse user response', err);
		})) as UserEndpointsResponse;

		if (!urls.endpoints || !urls.endpoints.length) {
			return redirect(301, '/onboarding');
		}

		return urls;
	};

	return { user: await fetchUser(), url: await fetchUrls(), err: null };
};
