import { z } from 'zod';
import type { Actions, PageServerLoad } from './$types';
import { fail, message, superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';
import { error, redirect } from '@sveltejs/kit';
import { durationStringToGoDurationString, Hour, toGoDurationString } from '$lib/utils';
import { getLocalTimeZone, parseDate } from '@internationalized/date';
import type { Category, CategoryTree, DurationString, UpdateTimeEntryInput } from '$lib/types';

const createTimeEntrySchema = z.object({
	categoryId: z.coerce.number(),
	date: z.string(),
	durationHours: z
		.union([z.coerce.number(), z.string().regex(/^(?:\d+|\d+t \d+m|0t)$/)])
		.default(0),
	description: z.string().optional()
});

const updateTimeEntrySchema = z.object({
	id: z.coerce.number(),
	durationHours: z
		.union([z.coerce.number(), z.string().regex(/^(?:\d+|\d+t \d+m|0t)$/)])
		.default(0),
	description: z.string().optional()
});

const deleteTimeEntrySchema = z.object({
	id: z.coerce.number()
});

export const load: PageServerLoad = async ({ locals, cookies, url, params }) => {
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
		if (daySummaryRes.status == 401) {
			cookies.delete('authToken', {path: '/'})
			redirect(303, `/auth/login?redirect=${url.pathname}`)
		}
		error(daySummaryRes.status, daySummaryRes.error);
	}

	const categoriesRes = await locals.apiService.getCategories(locals.authToken);
	if (!categoriesRes.ok) {
		if (categoriesRes.status == 401) {
			cookies.delete('authToken', {path: '/'})
			redirect(303, `/auth/login?redirect=${url.pathname}`)
		}
		error(categoriesRes.status, categoriesRes.error);
	}

	const categories = categoriesRes.data.flatMap(tree => flattenCategoryTree(tree));

	return {
		createForm,
		categories,
		daySummary: daySummaryRes.data
	};
};

const flattenCategoryTree = (tree: CategoryTree): Category[] => {
	const flattenChild = (
		tree: CategoryTree,
		parentTitles: string[],
		isParentFollowed: boolean
	): Category[] => {
		const children = tree.children.flatMap((child) =>
			flattenChild(child, [...parentTitles, tree.title], isParentFollowed || tree.isFollowed)
		);

		if (children.length > 0) {
			return children;
		} else if (tree.isFollowed || isParentFollowed) {
			return [
				{
					id: tree.id,
					title: tree.title,
					rootTitle: parentTitles.join(' - ')
				}
			];
		} else {
			return [];
		}
	};

	const children = tree.children.flatMap(child => flattenChild(child, [tree.title], tree.isFollowed))

	if (children.length > 0) {
		return children;
	} else if (tree.isFollowed) {
		return [
			{
				id: tree.id,
				title: tree.title,
				rootTitle: tree.title
			}
		];
	} else {
		return [];
	}
};

export const actions: Actions = {
	createTimeEntry: async ({ request, locals, cookies, url }) => {
		const form = await superValidate(request, zod(createTimeEntrySchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const duration =
			typeof form.data.durationHours == 'number'
				? toGoDurationString(form.data.durationHours * Hour)
				: durationStringToGoDurationString(form.data.durationHours as DurationString);

		const data = {
			date: form.data.date,
			duration,
			categoryId: form.data.categoryId,
			description: form.data.description
		};

		const res = await locals.apiService.createTimeEntry(data, locals.authToken);

		if (res.ok) {
			return message(form, 'Created Time Entry');
		}

		if (res.status == 401) {
			cookies.delete('authToken', {path: '/'})
			redirect(303, `/auth/login?redirect=${url.pathname}`)
		}

		return message(form, res.error, { status: 500 });
	},
	updateTimeEntry: async ({ request, locals, cookies, url }) => {
		const form = await superValidate(request, zod(updateTimeEntrySchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const duration =
			typeof form.data.durationHours == 'number'
				? toGoDurationString(form.data.durationHours * Hour)
				: durationStringToGoDurationString(form.data.durationHours as DurationString);

		const data: UpdateTimeEntryInput = {
			duration,
			description: form.data.description ?? ''
		};

		const res = await locals.apiService.updateTimeEntry(form.data.id, data, locals.authToken);

		if (res.ok) {
			return message(form, 'Updated Time Entry');
		}

		if (res.status == 401) {
			cookies.delete('authToken', {path: '/'})
			redirect(303, `/auth/login?redirect=${url.pathname}`)
		}

		return message(form, res.error, { status: 500 });
	},
	deleteTimeEntry: async ({ request, locals, cookies, url }) => {
		const form = await superValidate(request, zod(deleteTimeEntrySchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const entryId = form.data.id;

		const res = await locals.apiService.deleteTimeEntry(entryId, locals.authToken);

		if (res.ok) {
			return message(form, 'Deleted Time Entry');
		}

		if (res.status == 401) {
			cookies.delete('authToken', {path: '/'})
			redirect(303, `/auth/login?redirect=${url.pathname}`)
		}

		return message(form, res.error, { status: 500 });
	}
};
