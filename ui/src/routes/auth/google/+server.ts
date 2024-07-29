import { PUBLIC_BASE_URL } from '$env/static/public';
import { redirect } from '@sveltejs/kit';

export async function GET(): Promise<Response> {
	throw redirect(302, `${PUBLIC_BASE_URL}/auth/google`);
}
