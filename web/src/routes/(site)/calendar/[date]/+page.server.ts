import { z } from 'zod';
import type { Actions, PageServerLoad } from './$types';
import { fail, message, superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';
import { error } from '@sveltejs/kit';
import { Hour } from '$lib/utils';
import { getLocalTimeZone, parseDate } from '@internationalized/date';
import type { UpdateTimeEntryInput } from '$lib/types';

const createTimeEntrySchema = z.object({
	categoryId: z.coerce.number(),
	date: z.string(),
	durationHours: z.coerce.number(),
	description: z.string().optional()
});

const updateTimeEntrySchema = z.object({
	id: z.coerce.number(),
	durationHours: z.coerce.number(),
	description: z.string().optional()
});

const deleteTimeEntrySchema = z.object({
	id: z.coerce.number()
});

export const load: PageServerLoad = async ({ locals, params }) => {
	const defaultValues: z.infer<typeof createTimeEntrySchema> = {
		date: params.date,
		categoryId: -1,
		durationHours: 0
	};
	const createForm = await superValidate(defaultValues, zod(createTimeEntrySchema));

	const daySummaryRes = await locals.apiService.getSummaryForDate(
		parseDate(params.date).toDate(getLocalTimeZone()),
		locals.authToken
	);
	if (!daySummaryRes.ok) {
		error(500, daySummaryRes.error);
	}

	const categoryRes = await locals.apiService.getUserCategories(locals.authToken);
	if (!categoryRes.ok) {
		error(500, categoryRes.error);
	}

	return { createForm, categories: categoryRes.data, daySummary: daySummaryRes.data };
};

export const actions: Actions = {
	createTimeEntry: async ({ request, locals }) => {
		const form = await superValidate(request, zod(createTimeEntrySchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const data = {
			date: form.data.date,
			duration: Math.round(form.data.durationHours * Hour),
			categoryId: form.data.categoryId,
			description: form.data.description
		};

		const res = await locals.apiService.createTimeEntry(data, locals.authToken);

		if (res.ok) {
			return message(form, 'Created Time Entry');
		}

		return message(form, res.error, { status: 500 });
	},
	updateTimeEntry: async ({ request, locals }) => {
		const form = await superValidate(request, zod(updateTimeEntrySchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const data: UpdateTimeEntryInput = {
			duration: Math.round(form.data.durationHours * Hour),
			description: form.data.description ?? ''
		};

		const res = await locals.apiService.updateTimeEntry(form.data.id, data, locals.authToken);

		if (res.ok) {
			return message(form, 'Updated Time Entry');
		}

		return message(form, res.error, { status: 500 });
	},
	deleteTimeEntry: async ({ request, locals }) => {
		const form = await superValidate(request, zod(deleteTimeEntrySchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const entryId = form.data.id;

		const res = await locals.apiService.deleteTimeEntry(entryId, locals.authToken);

		if (res.ok) {
			return message(form, 'Deleted Time Entry');
		}

		return message(form, res.error, { status: 500 });
	}
};
