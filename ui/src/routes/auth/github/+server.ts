import { redirect } from "@sveltejs/kit";

export async function GET(): Promise<Response> {
	// TODO: Replace this URL with actual endpoint 
	return redirect(302, "http://localhost:3000/auth/github")
}
