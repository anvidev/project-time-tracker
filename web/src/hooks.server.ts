import { ApiServiceFactory } from '$lib/apiService';
import { redirect, type Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.apiService = ApiServiceFactory(event.fetch, 'http://localhost:9090');

	const authToken = event.cookies.get('authToken')
	if (!event.url.pathname.startsWith('/auth') && authToken == undefined) {
		redirect(303, '/auth/login')
	}

	event.locals.authToken = authToken ?? ""

	return await resolve(event);
};
