import { fail, superValidate } from 'sveltekit-superforms';
import { zod } from 'sveltekit-superforms/adapters';
import { z } from 'zod';
import type { Actions, PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';

const schema = z.object({
	email: z.string().email(),
	password: z.string().min(8).max(32)
});

export const load: PageServerLoad = async () => {
	const form = await superValidate(zod(schema));

	return { form };
};

export const actions: Actions = {
	default: async ({ request, locals, cookies }) => {
		const form = await superValidate(request, zod(schema));
		if (!form.valid) {
			return fail(400, { form });
		}

		const loginRes = await locals.apiService.logIn(form.data);
		if (!loginRes.ok) {
			return fail(500, { form, error: loginRes.error });
		}

		cookies.set('authToken', loginRes.data.token, { path: '/', maxAge: (7 * 24 * 60 * 60) });

		redirect(303, '/calendar');
	}
};
