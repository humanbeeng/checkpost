import { PUBLIC_BASE_URL } from '$env/static/public';
import type { User } from '@/types';
import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from '../$types';
import type { GenerateEndpointResponse, UserEndpointsResponse } from './types';

export const load: PageServerLoad = async ({ fetch, cookies }) => {
	const token = cookies.get('token');
	if (!token) {
		return redirect(301, '/');
	}

	// TODO: Better error handling
	const fetchUser = async () => {
		const res = await fetch(`${PUBLIC_BASE_URL}/user`).catch((err) => {
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

	const fetchUserEndpoints = async () => {
		const res = await fetch(`${PUBLIC_BASE_URL}/endpoint`).catch((err) => {
			throw error(500);
		});

		if (!res.ok) {
			throw error(res.status, { message: await res.text() });
		}

		const endpoints = (await res.json().catch((err) => {
			console.log('Unable to parse user endpoints response', err);
			throw error(500, { message: 'Something went wrong' });
		})) as UserEndpointsResponse;

		return endpoints;
	};

	const user = await fetchUser();
	const userEndpoints = await fetchUserEndpoints();

	if (user && userEndpoints && userEndpoints.endpoints) {
		const endpoint = userEndpoints.endpoints.at(0);
		if (endpoint) {
			return redirect(301, `/inspect/${endpoint.endpoint}`);
		} else {
			return redirect(301, `/onboarding`);
		}
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
			const res = await fetch(`${PUBLIC_BASE_URL}/endpoint/generate`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(req)
			});

			if (res.ok) {
				return { endpoint: (await res.json()) as GenerateEndpointResponse, err: null };
			} else {
				const text = await res.text();
				return fail(res.status, { err: { field: '', message: text } });
			}
		} catch (err) {
			console.error('Error', err);
			return fail(500, { err: { field: '', message: 'Something went wrong' } });
		}
	}
} satisfies Actions;
