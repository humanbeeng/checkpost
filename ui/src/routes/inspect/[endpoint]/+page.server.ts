import { PUBLIC_BASE_URL } from '$env/static/public';
import type { EndpointHistory, User, WebsocketTokenResponse } from '@/types';
import { error } from '@sveltejs/kit';

export const csr = true;

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

	const fetchEndpointHistory = async () => {
		console.log('Fetching URL history');
		const res = await fetch(`${PUBLIC_BASE_URL}/endpoint/history/${endpoint}`).catch((err) => {
			console.error('Unable to fetch endpoint request history', err);
			throw error(500);
		});

		if (!res.ok) {
			throw error(res.status, { message: await res.text() });
		}

		const endpointHistory = (await res.json().catch((err) => {
			console.error('Unable to parse endpoint history', err);
		})) as EndpointHistory;

		return endpointHistory;
	};

	const fetchWSToken = async () => {
		console.log('Fetching websocket token');

		const res = await fetch(`${PUBLIC_BASE_URL}/endpoint/${endpoint}/generate-token`).catch(
			(err) => {
				console.error('Unable to fetch websocket token', err);
				throw error(500);
			}
		);

		if (!res.ok) {
			console.log('Not ok');
			throw error(res.status, { message: await res.text() });
		}

		const token = (await res.json().catch((err) => {
			console.error('Unable to parse websocket token', err);
		})) as WebsocketTokenResponse;

		return token.token;
	};

	const user = await fetchUser();
	const endpointHistory = await fetchEndpointHistory();
	const token = await fetchWSToken();

	return {
		user,
		endpointHistory,
		token
	};
};
