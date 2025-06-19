import { z } from "zod";
import type { Actions, PageServerLoad } from "./$types";
import { fail, message, superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { error } from "@sveltejs/kit";
import { Hour } from "$lib/utils";
import { format } from "date-fns";

const schema = z.object({
	categoryId: z.coerce.number(),
	date: z.string(),
	durationHours: z.coerce.number(),
	description: z.string().optional(),
})

export const load: PageServerLoad = async ({ locals }) => {
	const defaultValues: z.infer<typeof schema> = {
		date: format(new Date(), "yyyy-MM-dd"),
		categoryId: -1,
		durationHours: 0,
	}
	const form = await superValidate(defaultValues, zod(schema));

	const daySummaryRes = await locals.apiService.getSummaryForDate(new Date(), locals.authToken)
	if (!daySummaryRes.ok) {
		error(500, daySummaryRes.error)
	}

	const categoryRes = await locals.apiService.getUserCategories(locals.authToken)
	if (!categoryRes.ok) {
		error(500, categoryRes.error)
	}

	return { form, categories: categoryRes.data, daySummary: daySummaryRes.data };
}

export const actions: Actions = {
	default: async ({ request, locals }) => {
		const form = await superValidate(request, zod(schema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const data = {
			date: form.data.date,
			duration: form.data.durationHours * Hour,
			categoryId: form.data.categoryId,
			description: form.data.description,
		}

		const res = await locals.apiService.createTimeEntry(data, locals.authToken)

		if (res.ok) {
			return message(form, "Created Time Entry")
		} 

		return message(form, res.error, {status: 500})
	},
}
