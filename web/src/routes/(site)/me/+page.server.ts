import { error } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { z } from 'zod';
import { fail, message, superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';
import type { ServiceResponse } from '$lib/apiService';
import type { CategoryTree } from '$lib/types';

const toggleFollowSchema = z.object({
	id: z.coerce.number(),
	isFollowed: z.coerce.boolean()
});

const createCategorySchema = z.object({
	parentId: z.coerce.number().nullable(),
	title: z.string()
});

export const load: PageServerLoad = async ({ locals }) => {
	const categoryTrees = await locals.apiService.getCategories(locals.authToken);

	if (!categoryTrees.ok) {
		error(500, categoryTrees.error);
	}

	return { categoryTrees: sortCategoryTrees(categoryTrees.data) };
};

const sortCategoryTrees = (trees: CategoryTree[]) => {
	const sortTree = (tree: CategoryTree): CategoryTree => ({
			...tree,
			children: tree.children.map(child => sortTree(child)).sort((a, b) => a.id - b.id)
		})
	
	return trees.map(tree => sortTree(tree)).sort((a, b) => a.id - b.id)
}

export const actions: Actions = {
	toggleFollow: async ({ request, locals }) => {
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
			return fail(500, { form });
		}

		return message(form, 'Toggled follow');
	},
	createCategory: async ({ request, locals }) => {
		const form = await superValidate(request, zod(createCategorySchema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const res = await locals.apiService.createCategory(form.data, locals.authToken)
		if (!res.ok) {
			return fail(500, { form, message: res.error })
		}

		return message(form, "Created category")
	}
};
