import { PUBLIC_BASE_URL } from '$env/static/public';
import { error, redirect } from '@sveltejs/kit';
import type { RequestEvent } from '../$types';

export async function GET({ url, fetch, cookies }: RequestEvent) {
	// TODO: Handle edge cases
	const code = url.searchParams.get('code');
	const endpoint = `${PUBLIC_BASE_URL}/auth/github/callback?code=${code}`;

	const res = await fetch(endpoint).catch((err) => {
		console.error('Unable to hit auth callback', err);
		return error(500, { message: 'Something went wrong while callback' });
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

		return redirect(302, '/onboarding');
	} else {
		return error(401);
	}
}
