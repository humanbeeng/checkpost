import { error, redirect } from "@sveltejs/kit";
import type { RequestEvent } from "../$types";
import { accessToken } from "../../../../stores/auth";


export async function GET({ url, fetch }: RequestEvent) {
	const code = url.searchParams.get('code')
	// TODO: Handle edge cases

	// TODO: Replace this url with actual endpoint
	const endpoint = `http://localhost:3000/auth/github/callback?code=${code}`

	// TODO: Handle error case

	const res = await fetch(endpoint)
	if (res.ok) {
		const response = await res.json()
		const token = response.token;

		accessToken.set(token)

		redirect(302, '/')
	} else {
		error(401)
	}
}
