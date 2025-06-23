import { format, parse } from 'date-fns';
import type { PageServerLoad } from './$types';
import { error } from '@sveltejs/kit';
import type { Month } from '$lib/types';
import { getLocalTimeZone, parseDate } from '@internationalized/date';

const monthMap: Record<Month, number> = {
	january: 1,
	february: 2,
	march: 3,
	april: 4,
	may: 5,
	june: 6,
	july: 7,
	august: 8,
	september: 9,
	october: 10,
	november: 11,
	december: 12
};

export const load: PageServerLoad = async ({ locals, url }) => {
	const searchDate = url.searchParams.get('date');
	const date = searchDate == undefined ? new Date() : parse(searchDate, 'yyyy-MM-dd', new Date());

	const monthStr = format(date, 'yyyy-MM');

	const res = await locals.apiService.getSummaryForMonth(monthStr, locals.authToken);
	if (!res.ok) {
		error(500, res.error);
	}

	const calendarRes = await locals.apiService.getCalendarYear(date.getFullYear());
	if (!calendarRes.ok) {
		error(500, calendarRes.error);
	}

	const calendar = {
		days: calendarRes.data.days.filter((day) => day.month == monthMap[res.data.month])
	};

	return { summary: res.data, calendar, date };
};
