import { ApiServiceFactory } from '$lib/apiService';
import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.apiService = ApiServiceFactory(event.fetch, 'http://localhost:9090');

	return await resolve(event);
};
