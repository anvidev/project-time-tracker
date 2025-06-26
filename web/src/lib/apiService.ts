import { format } from 'date-fns';
import type {
	Calendar,
	Category,
	CategoryTree,
	NewCategory,
	RegisterTimeEntryInput,
	Session,
	SummaryDay,
	SummaryDayDTO,
	SummaryMonth,
	SummaryMonthDTO,
	TimeEntry,
	UpdateTimeEntryInput
} from './types';
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

export type ApiService = {
	register: (data: { name: string; email: string; password: string }) => Promise<ServiceResponse>;
	logIn: (data: { email: string; password: string }) => Promise<ServiceResponse<Session>>;
	getUserCategories: (authToken: string) => Promise<ServiceResponse<Category[]>>;
	getCategories: (authToken: string) => Promise<ServiceResponse<CategoryTree[]>>;
	createCategory: (data: NewCategory, authToken: string) => Promise<ServiceResponse<Category>>;
	followCategory: (id: number, authToken: string) => Promise<ServiceResponse<null>>;
	unfollowCategory: (id: number, authToken: string) => Promise<ServiceResponse<null>>;
	createTimeEntry: (
		data: RegisterTimeEntryInput,
		authToken: string
	) => Promise<ServiceResponse<TimeEntry>>;
	updateTimeEntry: (
		id: number,
		data: UpdateTimeEntryInput,
		authToken: string
	) => Promise<ServiceResponse<TimeEntry>>;
	deleteTimeEntry: (id: number, authToken: string) => Promise<ServiceResponse<undefined>>;
	getSummaryForDate: (
		date: Date | string,
		authToken: string
	) => Promise<ServiceResponse<SummaryDay>>;
	getSummaryForMonth: (
		monthStr: string,
		authToken: string
	) => Promise<ServiceResponse<SummaryMonth>>;
	getCalendarYear: (year: number) => Promise<ServiceResponse<Calendar>>;
};

let apiServiceInstance: ApiService | undefined;

export type TApiServiceFactory = (fetch: FetchFn, baseUrl: string) => ApiService;

export const ApiServiceFactory: TApiServiceFactory = (fetch: FetchFn, baseUrl: string) => {
	if (apiServiceInstance == undefined) {
		apiServiceInstance = {
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
					}
				});

				if (res.ok) {
					return {
						ok: true,
						data: (await res.json()).categories
					};
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
					const body = await res.text();

					console.error(body);
					return {
						ok: false,
						error: body
					};
				}
			},

			getCategories: async function(authToken: string): Promise<ServiceResponse<CategoryTree[]>> {
				const res = await fetch(`${baseUrl}/v1/me/categories/all`, {
					headers: {
						Authorization: `Bearer ${authToken}`
					}
				});

				if (res.ok) {
					return {
						ok: true,
						data: (await res.json()).categories
					};
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
					const body = await res.text();

					console.error(body);
					return {
						ok: false,
						error: body
					};
				}
			},
			createCategory: async function(data: NewCategory, authToken: string): Promise<ServiceResponse<Category>> {
				const res = await fetch(`${baseUrl}/v1/me/categories`, {
					method: 'POST',
					headers: {
						Authorization: `Bearer ${authToken}`
					},
					body: JSON.stringify(data)
				});

				if (res.ok) {
					return {
						ok: true,
						data: (await res.json()).category
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

			followCategory: async function(
				id: number,
				authToken: string
			): Promise<ServiceResponse<null>> {
				const res = await fetch(`${baseUrl}/v1/me/categories/${id}/follow`, {
					method: 'put',
					headers: {
						Authorization: `Bearer ${authToken}`
					}
				});

				if (res.ok) {
					return {
						ok: true,
						data: null
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
			unfollowCategory: async function(
				id: number,
				authToken: string
			): Promise<ServiceResponse<null>> {
				const res = await fetch(`${baseUrl}/v1/me/categories/${id}/unfollow`, {
					method: 'put',
					headers: {
						Authorization: `Bearer ${authToken}`
					}
				});

				if (res.ok) {
					return {
						ok: true,
						data: null
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

			// TIME_ENTRIES
			createTimeEntry: async function(
				data: RegisterTimeEntryInput,
				authToken: string
			): Promise<ServiceResponse<TimeEntry>> {
				const res = await fetch(`${baseUrl}/v1/me/time_entries`, {
					method: 'POST',
					headers: {
						Authorization: `Bearer ${authToken}`
					},
					body: JSON.stringify(data)
				});

				if (res.ok) {
					return {
						ok: true,
						data: await res.json().then((json) => ({
							...json.timeEntry,
							duration: parseDuration(json.timeEntry.duration)
						}))
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
			updateTimeEntry: async function(
				id: number,
				data: UpdateTimeEntryInput,
				authToken: string
			): Promise<ServiceResponse<TimeEntry>> {
				const res = await fetch(`${baseUrl}/v1/me/time_entries/${id}`, {
					method: 'put',
					headers: {
						Authorization: `Bearer ${authToken}`
					},
					body: JSON.stringify(data)
				});

				if (res.ok) {
					return {
						ok: true,
						data: await res.json().then((json) => ({
							...json.timeEntry,
							duration: parseDuration(json.timeEntry.duration)
						}))
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
			deleteTimeEntry: async function(
				id: number,
				authToken: string
			): Promise<ServiceResponse<undefined>> {
				const res = await fetch(`${baseUrl}/v1/me/time_entries/${id}`, {
					method: 'DELETE',
					headers: {
						Authorization: `Bearer ${authToken}`
					}
				});

				if (res.ok) {
					return {
						ok: true,
						data: undefined
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
			getSummaryForDate: async function(
				date: Date | string,
				authToken: string
			): Promise<ServiceResponse<SummaryDay>> {
				let dateStr: string;
				if (date instanceof Date) {
					dateStr = format(date, 'yyyy-MM-dd');
				} else {
					dateStr = date;
				}

				const res = await fetch(`${baseUrl}/v1/me/time_entries/day/${dateStr}`, {
					headers: {
						Authorization: `Bearer ${authToken}`
					}
				});

				if (res.ok) {
					return {
						ok: true,
						data: await res
							.json()
							.then((json) => json.summary as SummaryDayDTO)
							.then((summary) => {
								return {
									date: summary.date,
									totalHours: parseDuration(summary.totalHours) ?? -1,
									maxHours: parseDuration(summary.maxHours) ?? 0,
									timeEntries: summary.timeEntries.map((e) => ({
										...e,
										duration: parseDuration(e.duration) ?? -1
									}))
								};
							})
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
			getSummaryForMonth: async function(
				monthStr: string,
				authToken: string
			): Promise<ServiceResponse<SummaryMonth>> {
				const res = await fetch(`${baseUrl}/v1/me/time_entries/month/${monthStr}`, {
					headers: {
						Authorization: `Bearer ${authToken}`
					}
				});

				if (res.ok) {
					return {
						ok: true,
						data: await res
							.json()
							.then((json) => json.summary as SummaryMonthDTO)
							.then((summary) => ({
								month: summary.month,
								totalHours: parseDuration(summary.totalHours) ?? -1,
								maxHours: parseDuration(summary.maxHours) ?? -1,
								days: summary.days.map((day) => {
									return {
										date: day.date,
										totalHours: parseDuration(day.totalHours) ?? -1,
										maxHours: parseDuration(day.maxHours) ?? 0,
										timeEntries: day.timeEntries.map((e) => ({
											...e,
											duration: parseDuration(e.duration) ?? -1
										}))
									};
								})
							}))
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
			getCalendarYear: async function(year: number): Promise<ServiceResponse<Calendar>> {
				const res = await fetch(`https://api.kalendarium.dk/MinimalCalendar/${year}`);

				if (res.ok) {
					return {
						ok: true,
						data: await res.json()
					};
				}

				return {
					ok: false,
					error: await res.text()
				};
			}
		};
	}

	return apiServiceInstance!;
};
