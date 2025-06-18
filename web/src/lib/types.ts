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

export interface TimeEntry {
	id: number;
	categoryId: number;
	category: string;
	userId: number;
	date: string;
	duration: Duration;
	description: string;
}

export type Duration = number
