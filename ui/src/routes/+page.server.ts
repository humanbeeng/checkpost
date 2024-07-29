import { PUBLIC_BASE_URL } from '$env/static/public';
import type { User } from '@/types';
import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import type { UserEndpointsResponse } from './onboarding/types';

type EndpointExistsResponse = {
	endpoint: string;
	exists: boolean;
	message: string;
};

export const load: PageServerLoad = async ({ cookies, fetch }) => {
	// TODO: Add error returns
	const cookie = cookies.get('token');
	if (cookie) {
		const fetchUser = async () => {
			const res = await fetch(`${PUBLIC_BASE_URL}/user`).catch((err) => {
				error(500);
			});

			if (!res.ok) {
				error(res.status, { message: await res.text() });
			}

			const user = (await res.json().catch((err) => {
				console.log('Unable to parse user response', err);
				error(500, { message: 'Something went wrong' });
			})) as User;

			return user;
		};

		const fetchUserEndpoints = async () => {
			const res = await fetch(`${PUBLIC_BASE_URL}/endpoint`).catch((err) => {
				error(500);
			});

			if (!res.ok) {
				error(res.status, { message: await res.text() });
			}

			const endpoints = (await res.json().catch((err) => {
				console.log('Unable to parse user endpoints response', err);
				error(500, { message: 'Something went wrong' });
			})) as UserEndpointsResponse;

			return endpoints;
		};

		const user = await fetchUser();
		const userEndpoints = await fetchUserEndpoints();

		if (user && userEndpoints && userEndpoints.endpoints) {
			const endpoint = userEndpoints.endpoints.at(0);
			if (endpoint) {
				throw redirect(301, `/inspect/${endpoint.endpoint}`);
			} else {
				throw redirect(301, `/onboarding`);
			}
		}
	}
};

export const actions = {
	check: async ({ request, fetch }) => {
		const data = await request.formData();
		const subdomain = data.get('subdomain');

		if (!subdomain) {
			return;
		}

		if (subdomain.toString().length < 4 || subdomain.toString().length > 10) {
			return fail(400, {
				error: 'Subdomain length should be between 4 and 10.',
			});
		}
		const res = await fetch(`${PUBLIC_BASE_URL}/endpoint/exists/${subdomain}`, {
			method: 'GET',
		});

		switch (res.status) {
			case 200: {
				return {
					data: (await res.json()) as EndpointExistsResponse,
				};
			}
			case 400: {
				return {
					error: await res.text(),
				};
			}
		}
	},
} satisfies Actions;
