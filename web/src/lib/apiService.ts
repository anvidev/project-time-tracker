import type { Category, RegisterTimeEntryInput, Session, TimeEntry } from './types';
import { parseDuration } from './utils';

type FetchFn = (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;

export type ErrorServiceResponse<E extends string> = {
	ok: false;
	error: E;
};
export type SuccessServiceResponse<T> = {
	ok: true;
	data: T;
};

export type ServiceResponse<T = unknown, E extends string = string> =
	| SuccessServiceResponse<T>
	| ErrorServiceResponse<E>;

export const ApiServiceFactory = (fetch: FetchFn, baseUrl: string) => {
	return {
		// AUTH FUNCTIONS
		register: async function(data: {
			name: string;
			email: string;
			password: string;
		}): Promise<ServiceResponse> {
			const res = await fetch(`${baseUrl}/v1/auth/register`, {
				method: 'POST',
				body: JSON.stringify(data)
			});

			if (res.ok) {
				return {
					ok: true,
					data: await res.json()
				};
			}

			const body: {
				error: string;
				code: string;
			} = await res.json();

			return {
				ok: false,
				error: `${body.code}: ${body.error}`
			};
		},
		logIn: async function(data: {
			email: string;
			password: string;
		}): Promise<ServiceResponse<Session>> {
			const res = await fetch(`${baseUrl}/v1/auth/login`, {
				method: 'POST',
				body: JSON.stringify(data)
			});

			if (res.ok) {
				return {
					ok: true,
					data: await res.json().then((data) => data.session)
				};
			}

			const body: {
				error: string;
				code: string;
			} = await res.json();

			return {
				ok: false,
				error: `${body.code}: ${body.error}`
			};
		},

		// CATEGORIES
		getUserCategories: async function(authToken: string): Promise<ServiceResponse<Category[]>> {
			const res = await fetch(`${baseUrl}/v1/me/categories`, {
				headers: {
					Authorization: `Bearer ${authToken}`
				},
			})

			if (res.ok) {
				return {
					ok: true,
					data: (await res.json()).categories
				}
			}

			if (res.headers.get('content-type')?.includes('application/json')) {
				const body: {
					error: string;
					code: string;
				} = await res.json();

				return {
					ok: false,
					error: `${body.code}: ${body.error}`
				};
			} else {
				const body = await res.text()

				console.error(body)
				return {
					ok: false,
					error: body,
				}
			}
		},

		// TIME_ENTRIES
		createTimeEntry: async function(data: RegisterTimeEntryInput, authToken: string): Promise<ServiceResponse<TimeEntry>> {
			const res = await fetch(`${baseUrl}/v1/me/time_entries`, {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${authToken}`
				},
				body: JSON.stringify(data),
			})

			if (res.ok) {
				return {
					ok: true,
					data: await res.json().then(json => ({
						...json.timeEntry,
						duration: parseDuration(json.timeEntry.duration),
					})),
				}
			}

			const body: {
				error: string;
				code: string;
			} = await res.json();

			return {
				ok: false,
				error: `${body.code}: ${body.error}`
			};
		}
	};
};
