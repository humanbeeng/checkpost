import { PUBLIC_SERVER_URL } from '$env/static/public';
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
	let cookie = cookies.get('token');
	if (cookie) {
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
	}
};

export const actions = {
	check: async ({ request, fetch }) => {
		const data = await request.formData();
		const subdomain = data.get('subdomain');

		if (!subdomain) {
			return;
		}

		console.log(subdomain.toString());

		if (subdomain.toString().length < 4 || subdomain.toString().length > 10) {
			console.log('failing');
			return fail(400, {
				error: 'Subdomain length should be between 4 and 10.'
			});
		}
		const res = await fetch(`http://api.checkpost.local:3000/url/exists/${subdomain}`, {
			method: 'GET'
		});

		switch (res.status) {
			case 200: {
				return {
					data: (await res.json()) as EndpointExistsResponse
				};
			}
			case 400: {
				return {
					error: await res.text()
				};
			}
		}
	}
} satisfies Actions;
