import { PUBLIC_BASE_URL } from '$env/static/public';
import type { UrlHistory, User } from '@/types';
import { error } from '@sveltejs/kit';

export const load = async ({ fetch, params }) => {
	const endpoint = params.endpoint;

	const fetchUser = async () => {
		console.log('Fetching user details');
		const res = await fetch(`${PUBLIC_BASE_URL}/user`).catch((err) => {
			console.error('Unable to fetch user details', err);
			throw error(500);
		});

		if (!res.ok) {
			throw error(res.status, { message: await res.text() });
		}

		const user = (await res.json().catch((err) => {
			console.error('Unable to parse user response', err);
			throw error(500, { message: 'Something went wrong' });
		})) as User;

		return user;
	};

	const fetchUrlHistory = async () => {
		console.log('Fetching URL history');
		const res = await fetch(`${PUBLIC_BASE_URL}/url/history/${endpoint}`).catch((err) => {
			console.error('Unable to fetch url request history', err);
			throw error(500);
		});

		if (!res.ok) {
			throw error(res.status, { message: await res.text() });
		}

		const urlHistory = (await res.json().catch((err) => {
			console.error('Unable to parse url history', err);
		})) as UrlHistory;

		return urlHistory;
	};

	const user = await fetchUser();
	const urlHistory = await fetchUrlHistory();

	return {
		user,
		urlHistory
	};
};
