import { type ApiService } from '$lib/apiService';
import type { User } from '$lib/types';
// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			apiService: ApiService;
			authToken: string;
			user: User
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export { };
