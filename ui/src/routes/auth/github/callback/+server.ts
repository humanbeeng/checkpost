import { PUBLIC_BASE_URL } from '$env/static/public';
import { error, redirect } from '@sveltejs/kit';
import type { RequestEvent } from '../$types';

export async function GET({ url, fetch, cookies }: RequestEvent) {
	// TODO: Handle edge cases
	const code = url.searchParams.get('code');
	const endpoint = `${PUBLIC_BASE_URL}/auth/github/callback?code=${code}`;

	const res = await fetch(endpoint).catch((err) => {
		console.error('Unable to hit auth callback', err);
		error(500, { message: 'Something went wrong while callback' });
	});

	if (res.ok) {
		console.log('Setting cookie token');
		const response = await res.json();

		cookies.set('token', response.token, {
			path: '/',
			// TODO: Fetch expiry from response
			httpOnly: true,
			secure: process.env.NODE_ENV === 'production',
			domain: process.env.NODE_ENV === 'production' ? `.${process.env.DOMAIN}` : "",
			maxAge: 1000 * 60 * 60 * 24 * 365 * 10, // 10 years
		});

		redirect(301, '/onboarding');
	} else {
		error(401);
	}
}
