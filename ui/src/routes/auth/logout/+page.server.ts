import { redirect } from '@sveltejs/kit';

export const load = async ({ cookies }) => {
	cookies.delete('token', { path: '/' });
	return redirect(302, '/');
};
