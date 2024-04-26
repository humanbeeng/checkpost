import { fail, redirect, type Actions } from '@sveltejs/kit';
import type { User, UserEndpointsResponse } from './types';

export const load = async ({ fetch }) => {
	const user: User = await fetch('http://localhost:3000/user', { credentials: 'include' })
		.then((res) => res.json())
		.catch((err) => {
			console.log('Unable to fetch user', err);
			return;
		});

	const urls: UserEndpointsResponse = await fetch('http://localhost:3000/url')
		.then((res) => res.json())
		.catch((err) => {
			console.log('Unable to fetch user endpoints', err);
			return;
		});

	if (user && urls.endpoints.length) {
		redirect(301, '/dashboard');
	}

	return user;
};

export const actions = {
	default: async ({ request, fetch }) => {
		console.log('form submit');
		const formData = await request.formData();
		const endpoint = formData.get('subdomain');
		if (!endpoint || endpoint == '') {
			return fail(400, { endpoint, missing: true });
		}
		console.log(endpoint);
		const req = {
			endpoint: endpoint
		};
		const res = await fetch('http://localhost:3000/url/generate', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(req)
		});

		console.log(req);

		if (res.ok) {
			console.log('ok');
			const data = await res.json();
			console.log(data);
		} else {
			const err = await res.text();
			console.log(err);
		}
	}
} satisfies Actions;
