import { ApiServiceFactory } from '$lib/apiService';
import { redirect, type Handle } from '@sveltejs/kit';
import { API_BASE_URL } from "$env/static/private";

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.apiService = ApiServiceFactory(event.fetch, API_BASE_URL);

	const authToken = event.cookies.get('authToken');
	if (!event.url.pathname.startsWith('/auth') && authToken == undefined) {
		redirect(303, '/auth/login');
	}
	if (event.url.pathname == '/') {
		redirect(307, '/calendar');
	}

	event.locals.authToken = authToken ?? '';

	return await resolve(event);
};
