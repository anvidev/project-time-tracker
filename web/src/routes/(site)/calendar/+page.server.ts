import { format, parse } from 'date-fns';
import type { PageServerLoad } from './$types';
import { error, redirect } from '@sveltejs/kit';
import { monthMap } from '$lib/types';

export const load: PageServerLoad = async ({ locals, cookies, url }) => {
	const searchDate = url.searchParams.get('date');
	const date = searchDate == undefined ? new Date() : parse(searchDate, 'yyyy-MM-dd', new Date());

	const monthStr = format(date, 'yyyy-MM');

	const res = await locals.apiService.getSummaryForMonth(monthStr, locals.authToken);
	if (!res.ok) {
		if (res.status == 401) {
			cookies.delete('authToken', {path: '/'})
			redirect(303, `/auth/login?redirect=${url.pathname}`)
		}
		error(res.status, res.error);
	}

	const calendarRes = await locals.apiService.getCalendarYear(date.getFullYear());
	if (!calendarRes.ok) {
		if (calendarRes.status == 401) {
			cookies.delete('authToken', {path: '/'})
			redirect(303, `/auth/login?redirect=${url.pathname}`)
		}
		error(calendarRes.status, calendarRes.error);
	}

	const calendar = {
		days: calendarRes.data.days.filter((day) => day.month == monthMap[res.data.month])
	};

	return { summary: res.data, calendar, date };
};
