import { error, redirect } from '@sveltejs/kit';
import type { RequestEvent } from '../$types';

export async function GET({ url, fetch, cookies }: RequestEvent) {
	const code = url.searchParams.get('code');
	console.log('Code', code);
	// TODO: Handle edge cases

	// TODO: Replace this url with actual endpoint
	const endpoint = `http://api.checkpost.local:3000/auth/github/callback?code=${code}`;

	// TODO: Handle error case

	const res = await fetch(endpoint).catch((err) => {
		console.log('Unable to hit auth callback', err);
		error(500, { message: 'Something went wrong while callback' });
	});

	if (res.ok) {
		const response = await res.json();
		// TODO: Increase security
		cookies.set('token', response.token, {
			path: '/',
			// TODO: Fetch expiry from response
			httpOnly: true,
			maxAge: 60 * 60 * 24 * 1000,
			secure: process.env.NODE_ENV === 'production'
		});

		redirect(302, '/onboarding');
	} else {
		error(401);
	}
}
