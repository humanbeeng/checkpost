import { fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

type EndpointExistsResponse = {
	endpoint: string;
	exists: boolean;
	message: string;
};

export const load: PageServerLoad = async ({ cookies }) => {
	let cookie = cookies.get('token');
	if (cookie) {
		redirect(301, '/dashboard');
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
