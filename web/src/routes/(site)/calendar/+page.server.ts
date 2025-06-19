import { format } from 'date-fns';
import type { PageServerLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	const monthStr = format(new Date(), 'yyyy-MM');

	const res = await locals.apiService.getSummaryForMonth(monthStr, locals.authToken);

	if (!res.ok) {
		error(500, res.error);
	}

	return { summary: res.data };
};
