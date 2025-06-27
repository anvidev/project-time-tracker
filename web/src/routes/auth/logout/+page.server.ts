import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = ({cookies}) => {
	cookies.delete('authToken', { path: '/' })
	redirect(307, '/auth/login')
}
