import { error } from "@sveltejs/kit"
import type { PageServerLoad } from "./$types"

export const load: PageServerLoad = async ({ locals }) => { 
	const categoryTrees = await locals.apiService.getCategories(locals.authToken)

	if (!categoryTrees.ok) {
		error(500, categoryTrees.error)
	}

	return { categoryTrees: categoryTrees.data }
}
