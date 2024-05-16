import { PUBLIC_BASE_URL } from '$env/static/public';
import { redirect } from '@sveltejs/kit';

export async function GET(): Promise<Response> {
	// TODO: Replace this URL with actual endpoint
	console.log(PUBLIC_BASE_URL);
	return redirect(302, `${PUBLIC_BASE_URL}/auth/github`);
}
