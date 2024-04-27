import { PUBLIC_SERVER_URL } from '$env/static/public';
import { redirect } from '@sveltejs/kit';

export async function GET(): Promise<Response> {
	// TODO: Replace this URL with actual endpoint
	console.log(PUBLIC_SERVER_URL);
	return redirect(302, `${PUBLIC_SERVER_URL}/auth/github`);
}
