import { error, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { z } from 'zod';
import { fail, message, superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';
import type { ServiceResponse } from '$lib/apiService';
import type { CategoryTree, WeekdayHours, WeekdayHoursDTO } from '$lib/types';
import { toGoDurationString, weekDayMap, weekDayToNum } from '$lib/utils';

const toggleFollowSchema = z.object({
	id: z.coerce.number(),
	isFollowed: z.coerce.boolean()
});

const createCategorySchema = z.object({
	parentId: z.coerce.number().nullable(),
	title: z.string()
});

const updateMaxHoursSchema = z.object({
	maxHours: z.array(
		z.object({
			weekday: z.enum(['Mandag', 'Tirsdag', 'Onsdag', 'Torsdag', 'Fredag', 'Lørdag', 'Søndag']),
			hours: z.coerce.number()
		})
	)
});

export const load: PageServerLoad = async ({ locals, cookies, url }) => {
	const hoursRes = await locals.apiService.getMaxHours(locals.authToken);
	if (!hoursRes.ok) {
		if (hoursRes.status == 401) {
			cookies.delete('authToken', { path: '/' });
			redirect(303, `/auth/login?redirect=${url.pathname}`);
		}

		error(hoursRes.status, hoursRes.error);
	}

	const categoryTrees = await locals.apiService.getCategories(locals.authToken);

	if (!categoryTrees.ok) {
		if (categoryTrees.status == 401) {
			cookies.delete('authToken', { path: '/' });
			redirect(303, `/auth/login?redirect=${url.pathname}`);
		}

		error(categoryTrees.status, categoryTrees.error);
	}

	return {
		categoryTrees: sortCategoryTrees(categoryTrees.data),
		maxHours: sortMaxHours(hoursRes.data)
	};
};

const sortMaxHours = (hours: WeekdayHours[]) => {
	return hours.sort((a, b) => weekDayMap[a.weekday] - weekDayMap[b.weekday]);
};

const sortCategoryTrees = (trees: CategoryTree[]) => {
	const sortTree = (tree: CategoryTree): CategoryTree => ({
		...tree,
		children: tree.children.map((child) => sortTree(child)).sort((a, b) => a.id - b.id)
	});

	return trees.map((tree) => sortTree(tree)).sort((a, b) => a.id - b.id);
};

export const actions: Actions = {
	toggleFollow: async ({ request, locals, cookies, url }) => {
		const form = await superValidate(request, zod(toggleFollowSchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		let res: ServiceResponse<null>;
		if (form.data.isFollowed) {
			res = await locals.apiService.unfollowCategory(form.data.id, locals.authToken);
		} else {
			res = await locals.apiService.followCategory(form.data.id, locals.authToken);
		}

		if (!res.ok) {
			if (res.status == 401) {
				cookies.delete('authToken', { path: '/' });
				redirect(303, `/auth/login?redirect=${url.pathname}`);
			}
			return fail(res.status, { form });
		}

		return message(form, 'Toggled follow');
	},
	createCategory: async ({ request, locals, cookies, url }) => {
		const form = await superValidate(request, zod(createCategorySchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const res = await locals.apiService.createCategory(form.data, locals.authToken);
		if (!res.ok) {
			if (res.status == 401) {
				cookies.delete('authToken', { path: '/' });
				redirect(303, `/auth/login?redirect=${url.pathname}`);
			}
			return fail(res.status, { form });
		}

		return message(form, 'Created category');
	},
	updateMaxHours: async ({ request, locals, cookies, url }) => {
		const form = await superValidate(request, zod(updateMaxHoursSchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const data: WeekdayHoursDTO[] = form.data.maxHours.map((h) => ({
			weekday: weekDayToNum(h.weekday, { sundayFirst: true }),
			hours: toGoDurationString(h.hours)
		}));

		const res = await locals.apiService.updateMaxHours(data, locals.authToken);
		if (!res.ok) {
			if (res.status == 401) {
				cookies.delete('authToken', { path: '/' });
				redirect(303, `/auth/login?redirect=${url.pathname}`);
			}
			return fail(res.status, { form });
		}

		return message(form, 'Max Hours updated');
	}
};
