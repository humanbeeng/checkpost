// See https://kit.svelte.dev/docs/types#app

// for information about these interfaces
declare global {
	namespace App {
		interface Error {
			code: number;
			message: Message;
		}
		// interface Locals {
		// }
		// interface PageData {
		// 	user: User | null;
		// 	urls: UserEndpointsResponse | null;
		// }
		interface ActionData {
			res: any;
			err: {
				code: number;
				message: string;
			};
		}
		// interface PageState {}
		// interface Platform {}
	}
}
