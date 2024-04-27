import { redirect } from '@sveltejs/kit';

export const load = async ({ cookies }) => {
	console.log('Logout called');
	cookies.delete('token', { path: '/' });
	return redirect(302, '/');
};
