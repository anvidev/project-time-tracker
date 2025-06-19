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

export interface SummaryMonth {
	month: string;
	totalHours: Duration;
	maxHours: Duration;
	days: SummaryDay[];
}

export interface SummaryMonthDTO {
	month: string;
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
