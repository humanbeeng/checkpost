import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from '../../inspect/[endpoint]/$types';

export const load: PageServerLoad = async ({ cookies }) => {
	console.log('Logout called');
	cookies.delete('token', { path: '/' });
	return redirect(302, '/');
};
