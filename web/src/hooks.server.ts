import { ApiServiceFactory } from '$lib/apiService';
import { redirect, type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.apiService = ApiServiceFactory(event.fetch, 'http://localhost:9090');

	if (!event.url.pathname.startsWith('/auth') && event.cookies.get('authToken') == undefined) {
		redirect(303, '/auth/login')
	}

	return await resolve(event);
};
