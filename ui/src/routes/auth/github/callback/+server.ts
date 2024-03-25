import { error, redirect } from "@sveltejs/kit";
import type { RequestEvent } from "../$types";


export async function GET({ url, fetch, cookies }: RequestEvent) {
	const code = url.searchParams.get('code')
	console.log("Code", code)
	// TODO: Handle edge cases

	// TODO: Replace this url with actual endpoint
	const endpoint = `http://localhost:3000/auth/github/callback?code=${code}`

	// TODO: Handle error case

	const res = await fetch(endpoint)

	if (res.ok) {
		const response = await res.json()

		// TODO: Increase security
		cookies.set('token', response.token, {
			path: '/',
			sameSite: 'strict',
			maxAge: 60 * 60 * 24 * 1000
		});

		redirect(302, '/dashboard')
	} else {
		error(401)
	}
}
