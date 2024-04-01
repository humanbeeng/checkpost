import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "../$types";

export const load: PageServerLoad = async ({ request, fetch }) => {
	const res = await fetch('http://localhost:3000/admin/dashboard')

	if (res.ok) {
		const data = await res.json()
		console.log(data)
	} else {
		error(400)
	}

}
