import type { Session } from './types';

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
		fetch,
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
		}
	};
};
