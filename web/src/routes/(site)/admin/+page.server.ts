import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, cookies, url }) => {
	const searchParams = url.searchParams;
	console.log('url', searchParams.toString());
	const entries = await locals.apiService.getAdminEntries(locals.authToken, searchParams);
	if (!entries.ok) {
		if (entries.status == 401) {
			cookies.delete('authToken', { path: '/' });
			redirect(303, `/auth/login?redirect=${url.pathname}`);
		}
		error(entries.status, entries.error);
	}

	const categoriesRes = await locals.apiService.getAllCategories(locals.authToken);
	if (!categoriesRes.ok) {
		if (categoriesRes.status == 401) {
			cookies.delete('authToken', { path: '/' });
			redirect(303, `/auth/login?redirect=${url.pathname}`);
		}
		error(categoriesRes.status, categoriesRes.error);
	}

	const usersRes = await locals.apiService.getAllUsers(locals.authToken);
	if (!usersRes.ok) {
		if (usersRes.status == 401) {
			cookies.delete('authToken', { path: '/' });
			redirect(303, `/auth/login?redirect=${url.pathname}`);
		}
		error(usersRes.status, usersRes.error);
	}

	return {
		entries: entries.data,
		categories: categoriesRes.data,
		users: usersRes.data
	};
};
