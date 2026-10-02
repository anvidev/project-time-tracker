import { format } from 'date-fns';
import type {
	AdminEntry,
	Calendar,
	Category,
	CategoryTree,
	LiveTimer,
	LiveTimerDTO,
	NewCategory,
	RegisterTimeEntryInput,
	SaveLiveTimerInput,
	Session,
	StartLiveTimerInput,
	SummaryDay,
	SummaryDayDTO,
	SummaryMonth,
	SummaryMonthDTO,
	TimeEntry,
	UpdateLiveTimerInput,
	UpdateTimeEntryInput,
	User,
	WeekdayHours,
	WeekdayHoursDTO
} from './types';
import { dayNumToWeekDay, parseDuration } from './utils';

type FetchFn = (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;

export type ErrorServiceResponse<E extends string> = {
	ok: false;
	status: number;
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
	getMaxHours: (authToken: string) => Promise<ServiceResponse<WeekdayHours[]>>;
	updateMaxHours: (data: WeekdayHoursDTO[], authToken: string) => Promise<ServiceResponse<null>>;
	getUserProfile: (authToken: string) => Promise<ServiceResponse<User>>;
	getAdminEntries: (
		authToken: string,
		searchParams: URLSearchParams
	) => Promise<ServiceResponse<AdminEntry>>;
	getAllCategories: (authToken: string) => Promise<ServiceResponse<Category[]>>;
	getAllUsers: (authToken: string) => Promise<ServiceResponse<User[]>>;
	getTimer: (authToken: string) => Promise<ServiceResponse<LiveTimer | null>>;
	startTimer: (data: StartLiveTimerInput, authToken: string) => Promise<ServiceResponse<LiveTimer>>;
	pauseTimer: (authToken: string) => Promise<ServiceResponse<LiveTimer>>;
	resumeTimer: (authToken: string) => Promise<ServiceResponse<LiveTimer>>;
	updateTimer: (
		data: UpdateLiveTimerInput,
		authToken: string
	) => Promise<ServiceResponse<LiveTimer>>;
	saveTimer: (data: SaveLiveTimerInput, authToken: string) => Promise<ServiceResponse<TimeEntry>>;
	discardTimer: (authToken: string) => Promise<ServiceResponse<undefined>>;
};

let apiServiceInstance: ApiService | undefined;

export type TApiServiceFactory = (fetch: FetchFn, baseUrl: string) => ApiService;

const parseLiveTimer = (dto: LiveTimerDTO): LiveTimer => ({
	...dto,
	elapsed: parseDuration(dto.elapsed) ?? 0
});

const timerError = async (res: Response): Promise<ErrorServiceResponse<string>> => {
	const body: { error: string; code: string } = await res.json();

	return {
		ok: false,
		status: res.status,
		error: `${body.code}: ${body.error}`
	};
};

export const ApiServiceFactory: TApiServiceFactory = (fetch: FetchFn, baseUrl: string) => {
	const timerRequest = async (
		method: string,
		path: string,
		authToken: string,
		data?: unknown
	): Promise<ServiceResponse<LiveTimer>> => {
		const res = await fetch(`${baseUrl}/v1/me/timer${path}`, {
			method,
			headers: { Authorization: `Bearer ${authToken}` },
			body: data === undefined ? undefined : JSON.stringify(data)
		});

		if (res.ok) {
			return { ok: true, data: parseLiveTimer((await res.json()).timer) };
		}

		return await timerError(res);
	};

	if (apiServiceInstance == undefined) {
		apiServiceInstance = {
			// AUTH FUNCTIONS
			register: async function (data: {
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
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},
			logIn: async function (data: {
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
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},

			// CATEGORIES
			getUserCategories: async function (authToken: string): Promise<ServiceResponse<Category[]>> {
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
						status: res.status,
						error: `${body.code}: ${body.error}`
					};
				} else {
					const body = await res.text();

					console.error(body);
					return {
						ok: false,
						status: res.status,
						error: body
					};
				}
			},

			getCategories: async function (authToken: string): Promise<ServiceResponse<CategoryTree[]>> {
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
						status: res.status,
						error: `${body.code}: ${body.error}`
					};
				} else {
					const body = await res.text();

					console.error(body);
					return {
						ok: false,
						status: res.status,
						error: body
					};
				}
			},
			createCategory: async function (
				data: NewCategory,
				authToken: string
			): Promise<ServiceResponse<Category>> {
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
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},

			followCategory: async function (
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
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},
			unfollowCategory: async function (
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
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},

			// TIME_ENTRIES
			createTimeEntry: async function (
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
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},
			updateTimeEntry: async function (
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
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},
			deleteTimeEntry: async function (
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
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},
			getSummaryForDate: async function (
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
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},
			getSummaryForMonth: async function (
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
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},
			getCalendarYear: async function (year: number): Promise<ServiceResponse<Calendar>> {
				const res = await fetch(`https://api.kalendarium.dk/MinimalCalendar/${year}`);

				if (res.ok) {
					return {
						ok: true,
						data: await res.json()
					};
				}

				return {
					ok: false,
					status: res.status,
					error: await res.text()
				};
			},
			getMaxHours: async function (authToken: string): Promise<ServiceResponse<WeekdayHours[]>> {
				const res = await fetch(`${baseUrl}/v1/me/hours`, {
					headers: {
						Authorization: `Bearer ${authToken}`
					}
				});

				if (res.ok) {
					return {
						ok: true,
						data: await res
							.json()
							.then((json) => json.hours as WeekdayHoursDTO[])
							.then((hours) =>
								hours.map((hour) => ({
									weekday: dayNumToWeekDay(hour.weekday, { sundayFirst: true }),
									hours: parseDuration(hour.hours) ?? 0
								}))
							)
					};
				}

				const body: {
					error: string;
					code: string;
				} = await res.json();

				return {
					ok: false,
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},
			updateMaxHours: async function (
				data: WeekdayHoursDTO[],
				authToken: string
			): Promise<ServiceResponse<null>> {
				const res = await fetch(`${baseUrl}/v1/me/hours`, {
					method: 'put',
					headers: {
						Authorization: `Bearer ${authToken}`
					},
					body: JSON.stringify({ hours: data })
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
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},
			getUserProfile: async function (authToken: string): Promise<ServiceResponse<User>> {
				const res = await fetch(`${baseUrl}/v1/me/profile`, {
					headers: {
						Authorization: `Bearer ${authToken}`
					}
				});

				if (res.ok) {
					return {
						ok: true,
						data: await res.json().then((json) => json.user as User)
					};
				}

				const body: {
					error: string;
					code: string;
				} = await res.json();

				return {
					ok: false,
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},
			getAdminEntries: async function (
				authToken: string,
				searchParams: URLSearchParams
			): Promise<ServiceResponse<AdminEntry>> {
				const res = await fetch(`${baseUrl}/v1/admin/time_entries?${searchParams.toString()}`, {
					headers: {
						Authorization: `Bearer ${authToken}`
					}
				});

				if (res.ok) {
					return {
						ok: true,
						data: (await res.json().then((json) => ({
							...json,
							timeSpent: parseDuration(json.timeSpent),
							entries: json.entries.map((entry: any) => ({
								...entry,
								duration: parseDuration(entry.duration)
							}))
						}))) as AdminEntry
					};
				}

				const body: {
					error: string;
					code: string;
				} = await res.json();

				return {
					ok: false,
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},
			getAllCategories: async function (authToken: string): Promise<ServiceResponse<Category[]>> {
				const res = await fetch(`${baseUrl}/v1/admin/categories`, {
					headers: {
						Authorization: `Bearer ${authToken}`
					}
				});

				if (res.ok) {
					return {
						ok: true,
						data: (await res.json().then((json) => json.categories)) as Category[]
					};
				}

				const body: {
					error: string;
					code: string;
				} = await res.json();

				return {
					ok: false,
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
			},
			getAllUsers: async function (authToken: string): Promise<ServiceResponse<User[]>> {
				const res = await fetch(`${baseUrl}/v1/admin/users`, {
					headers: {
						Authorization: `Bearer ${authToken}`
					}
				});

				if (res.ok) {
					return {
						ok: true,
						data: (await res.json().then((json) => json.users)) as User[]
					};
				}

				const body: {
					error: string;
					code: string;
				} = await res.json();

				return {
					ok: false,
					status: res.status,
					error: `${body.code}: ${body.error}`
				};
				},

				// LIVE TIMER
				getTimer: async function (authToken: string): Promise<ServiceResponse<LiveTimer | null>> {
					const res = await fetch(`${baseUrl}/v1/me/timer`, {
						headers: { Authorization: `Bearer ${authToken}` }
					});

					if (res.ok) {
						const json = await res.json();
						return { ok: true, data: json.timer ? parseLiveTimer(json.timer) : null };
					}

					return await timerError(res);
				},
				startTimer: async function (
					data: StartLiveTimerInput,
					authToken: string
				): Promise<ServiceResponse<LiveTimer>> {
					return await timerRequest('POST', '/start', authToken, data);
				},
				pauseTimer: async function (authToken: string): Promise<ServiceResponse<LiveTimer>> {
					return await timerRequest('POST', '/pause', authToken);
				},
				resumeTimer: async function (authToken: string): Promise<ServiceResponse<LiveTimer>> {
					return await timerRequest('POST', '/resume', authToken);
				},
				updateTimer: async function (
					data: UpdateLiveTimerInput,
					authToken: string
				): Promise<ServiceResponse<LiveTimer>> {
					return await timerRequest('PUT', '', authToken, data);
				},
				saveTimer: async function (
					data: SaveLiveTimerInput,
					authToken: string
				): Promise<ServiceResponse<TimeEntry>> {
					const res = await fetch(`${baseUrl}/v1/me/timer/save`, {
						method: 'POST',
						headers: { Authorization: `Bearer ${authToken}` },
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

					return await timerError(res);
				},
				discardTimer: async function (authToken: string): Promise<ServiceResponse<undefined>> {
					const res = await fetch(`${baseUrl}/v1/me/timer`, {
						method: 'DELETE',
						headers: { Authorization: `Bearer ${authToken}` }
					});

					if (res.ok) {
						return { ok: true, data: undefined };
					}

					return await timerError(res);
				}
		};
	}

	return apiServiceInstance!;
};
