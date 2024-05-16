import { PUBLIC_SERVER_URL } from '$env/static/public';
import type { User } from '@/types';
import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from '../$types';
import type { GenerateUrlResponse, UserEndpointsResponse } from './types';

export const load: PageServerLoad = async ({ fetch, cookies }) => {
	const token = cookies.get('token');
	if (!token) {
		return redirect(301, '/');
	}
	// TODO: Better error handling

	const fetchUser = async () => {
		const res = await fetch(`${PUBLIC_SERVER_URL}/user`).catch((err) => {
			throw error(500);
		});

		if (!res.ok) {
			throw error(res.status, { message: await res.text() });
		}

		const user = (await res.json().catch((err) => {
			console.log('Unable to parse user response', err);
			throw error(500, { message: 'Something went wrong' });
		})) as User;

		return user;
	};

	const fetchUserUrls = async () => {
		const res = await fetch(`${PUBLIC_SERVER_URL}/url`).catch((err) => {
			throw error(500);
		});

		if (!res.ok) {
			throw error(res.status, { message: await res.text() });
		}

		const urls = (await res.json().catch((err) => {
			console.log('Unable to parse user urls response', err);
			throw error(500, { message: 'Something went wrong' });
		})) as UserEndpointsResponse;

		return urls;
	};

	const user = await fetchUser();
	const urls = await fetchUserUrls();

	if (user && urls && urls.endpoints) {
		return redirect(301, '/waitlist');
	}
	return {
		user
	};
};

export const actions = {
	default: async ({ request, fetch }) => {
		const formData = await request.formData();
		const endpoint = formData.get('subdomain');
		if (!endpoint || endpoint == '') {
			return fail(400, { err: { field: 'subdomain', message: 'Empty subdomain' } });
		}
		const req = {
			endpoint: endpoint
		};
		try {
			const res = await fetch(`${PUBLIC_SERVER_URL}/url/generate`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(req)
			});

			if (res.ok) {
				return { url: (await res.json()) as GenerateUrlResponse, err: null };
			} else {
				const text = await res.text();
				return fail(res.status, { err: { field: '', message: text } });
			}
		} catch (err) {
			console.log('Error', err);
			return fail(500, { err: { field: '', message: 'Something went wrong' } });
		}
	}
} satisfies Actions;
