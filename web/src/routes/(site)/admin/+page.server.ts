import { error, redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, cookies, url }) => {
	const entries = await locals.apiService.getAdminEntries(locals.authToken);

	if (!entries.ok) {
		if (entries.status == 401) {
			cookies.delete('authToken', { path: '/' });
			redirect(303, `/auth/login?redirect=${url.pathname}`);
		}
		error(entries.status, entries.error);
	}

	return {
		entries: entries.data
	};
};
