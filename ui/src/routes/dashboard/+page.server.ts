import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from '../$types';

export const load: PageServerLoad = async ({ cookies }) => {
	console.log(cookies.get('token'));
	if (!cookies.get('token')) {
		redirect(301, '/');
	}
};
