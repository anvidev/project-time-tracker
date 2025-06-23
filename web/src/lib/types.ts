export interface Session {
	token: string;
	userId: number;
	expiresAt: string;
	createdAt: string;
	updatedAt: string;
}

export interface Category {
	id: number;
	title: string;
	rootTitle: string;
}

export interface RegisterTimeEntryInput {
	categoryId: number;
	date: string;
	duration: Duration;
	description?: string;
}

export interface SummaryDay {
	date: string;
	totalHours: Duration;
	maxHours: Duration;
	timeEntries: TimeEntry[];
}

export interface SummaryDayDTO {
	date: string;
	totalHours: string;
	maxHours: string;
	timeEntries: TimeEntryDTO[];
}

export const monthMap: Record<Month, number> = {
	january: 1,
	february: 2,
	march: 3,
	april: 4,
	may: 5,
	june: 6,
	july: 7,
	august: 8,
	september: 9,
	october: 10,
	november: 11,
	december: 12
}

export const months: Month[] = Object.keys(monthMap) as Month[]

export type Month =
	| 'january'
	| 'february'
	| 'march'
	| 'april'
	| 'may'
	| 'june'
	| 'july'
	| 'august'
	| 'september'
	| 'october'
	| 'november'
	| 'december';

export interface SummaryMonth {
	month: Month;
	totalHours: Duration;
	maxHours: Duration;
	days: SummaryDay[];
}

export interface SummaryMonthDTO {
	month: Month;
	totalHours: string;
	maxHours: string;
	days: SummaryDayDTO[];
}

export interface TimeEntry {
	id: number;
	categoryId: number;
	category: string;
	userId: number;
	date: string;
	duration: Duration;
	description: string;
}

export interface TimeEntryDTO {
	id: number;
	categoryId: number;
	category: string;
	userId: number;
	date: string;
	duration: string;
	description: string;
}

export type Duration = number;

export interface Calendar {
	days: Day[];
}

export interface Day {
	year: number;
	month: number;
	day: number;
	julianDay: number;
	date: string;
	formattedDate: string;
	dayInYear: number;
	dayName: string;
	altNames: string;
	wiki_link: string;
	weekday: string;
	weekNumber: string;
	holliday: boolean;
	moonSymbol: string;
	events: Event[];
}

export interface Event {
	id: number;
	danishShort: string;
	danishLong: string;
	latinShort: string;
	latinLong: string;
	merkedage: boolean;
	minimal: boolean;
	kirkeaar: boolean;
	holliday: boolean;
	liturgiskFarve: string;
	vises: string;
	visesHellig: string;
	definition: string;
	wikiLink: string;
}

export type WeekDay = 'Mandag' | 'Tirsdag' | 'Onsdag' | 'Torsdag' | 'Fredag' | 'Lørdag' | 'Søndag';
