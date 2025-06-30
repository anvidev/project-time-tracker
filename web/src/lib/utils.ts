import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';
import type { Duration, DurationString, GoDurationString, WeekDay } from './types';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, 'child'> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any } ? Omit<T, 'children'> : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & { ref?: U | null };

export const maxFractionDigits = (num: number, digits: number) => {
	if (num % 1 == 0) {
		return num;
	}

	const multiplier = Math.pow(10, digits);
	return Math.round(num * multiplier) / multiplier;
};

export const weekDayMap: Record<WeekDay, number> = {
	Mandag: 0,
	Tirsdag: 1,
	Onsdag: 2,
	Torsdag: 3,
	Fredag: 4,
	Lørdag: 5,
	Søndag: 6
};

export const numToWeekDayMap: Record<number, WeekDay> = {
	0: "Mandag",
	1: "Tirsdag",
	2: "Onsdag",
	3: "Torsdag",
	4: "Fredag",
	5: "Lørdag",
	6: "Søndag"
};

export const dayNumToWeekDay: (day: number, options?: {
	sundayFirst?: boolean
}) => WeekDay = (day, options = {}) => { 
	if (options.sundayFirst) {
		const moveDay = (day: number) => {
			if (day == 0) {
				return 6
			}
			return day-1
		}

		day = moveDay(day)
	}

	return numToWeekDayMap[day]
};

export const weekDayToNum: (day: WeekDay, options?: {
	sundayFirst?: boolean
}) => number = (day, options = {}) => { 
	let num = weekDayMap[day]

	if (options.sundayFirst) {
		const moveDay = (day: number) => {
			if (day == 6) {
				return 0
			}
			return day+1
		}

		num = moveDay(num)
	}

	return num
};

export const Nanosecond: Duration = 1;
export const Microsecond: Duration = 1000 * Nanosecond;
export const Millisecond: Duration = 1000 * Microsecond;
export const Second: Duration = 1000 * Millisecond;
export const Minute: Duration = 60 * Second;
export const Hour: Duration = 60 * Minute;

const unitMap = new Map<string, number>([
	['ns', Nanosecond],
	['us', Microsecond],
	['ms', Millisecond],
	['s', Second],
	['m', Minute],
	['h', Hour]
]);

export function toDurationString(d: Duration): DurationString {
	if (d == 0) return '0t';

	const hours = Math.floor(d / Hour);
	const minutes = Math.floor((d - hours * Hour) / Minute);

	return `${hours}t ${minutes}m`;
}

export function durationStringToGoDurationString(str: DurationString): GoDurationString {
	if (str == '0t') {
		return '0s';
	}

	const replaced = str.replace(' ', '').replace('t', 'h');

	return `${replaced as `${number}h${number}m`}0s`;
}

export function toGoDurationString(d: Duration): GoDurationString {
	if (d == 0) return '0s';

	const hours = Math.floor(d / Hour);
	const minutes = Math.floor((d - hours * Hour) / Minute);
	const seconds = Math.floor((d - (hours * Hour + minutes * Minute)) / Second);

	return `${hours}h${minutes}m${seconds}s`;
}

// This is a port of golangs time.ParseDuration.
// Dont ask me how this works, since i just ported it.
export function parseDuration(s: string): Duration | undefined {
	const orig = s;
	let d: Duration = 0;
	let neg = false;

	if (s != '') {
		const c = s[0];
		if (c == '-' || c == '+') {
			neg = c == '-';
			s = s.slice(1);
		}
	}

	if (s == '0') {
		return 0;
	}
	if (s == '') {
		console.error(`parseDuration: invalid duration '${orig}'. Empty string not allowed.`);
		return undefined;
	}

	while (s != '') {
		let v,
			f: number = 0;
		let scale: number = 1; // value = v + f/scale

		// The next character must be [0-9.]
		if (
			!(
				s[0] == '.' ||
				('0'.charCodeAt(0) <= s[0].charCodeAt(0) && s[0].charCodeAt(0) <= '9'.charCodeAt(0))
			)
		) {
			const pos = orig.length - s.length;
			console.error(
				`parseDuration: invalid duration '${orig}'. Expected number, got '${orig[pos]}'`
			);
			return undefined;
		}
		// Consume [0-9]*
		const pl = s.length;
		const leading = leadingInt(s);
		if (leading == undefined) {
			const pos = orig.length - s.length;
			console.error(`parseDuration: invalid duration '${orig}'. Invalid leading at ${pos}`);
			return undefined;
		}
		v = leading.x;
		s = leading.rem;
		const pre = pl != s.length; // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		let post = false;
		if (s != '' && s[0] == '.') {
			s = s.slice(1);
			const pl = s.length;
			const fract = leadingFraction(s);
			f = fract.x;
			scale = fract.scale;
			s = fract.rem;
			post = pl != s.length;
		}
		if (!pre && !post) {
			const pos = orig.length - s.length;
			console.error(`parseDuration: invalid duration '${orig}'. No pre/post at ${pos}`);
			return undefined;
		}

		// Consume unit.
		let i = 0;
		for (; i < s.length; i++) {
			const c = s[i];
			if (c == '.' || ('0' <= c && c <= '9')) {
				break;
			}
		}
		if (i == 0) {
			const pos = orig.length - s.length;
			console.error(`parseDuration: missing unit in duration '${orig}' at ${pos}`);
			return undefined;
		}
		const u = s.slice(0, i);
		s = s.slice(i);
		const unit = unitMap.get(u);
		if (unit == undefined) {
			const pos = orig.length - s.length;
			console.error(`parseDuration: unknown unit '${u}' in duration '${orig}' at ${pos}`);
			return undefined;
		}
		// if (v > 1<<63/unit) {
		// 	const pos = orig.length - s.length
		// 	console.error(`parseDuration: invalid duration '${orig}' at ${pos}`)
		// 	return undefined
		// }
		v *= unit;
		if (f > 0) {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += f * (unit / scale);
			// if (v > 1<<31) {
			// 	const pos = orig.length - s.length
			// 	console.error(`parseDuration: invalid duration '${orig}' at ${pos} (v > 1<<31)`)
			// 	return undefined
			// }
		}
		d += v;
		// if (d > 1<<31) {
		// 	const pos = orig.length - s.length
		// 	console.error(`parseDuration: invalid duration '${orig}' at ${pos} (d > 1<<31)`)
		// 	return undefined
		// }
	}

	return d;
}

// leadingInt consumes the leading [0-9]* from s.
function leadingInt(s: string): { x: number; rem: string } | undefined {
	let x: number = 0;

	let i = 0;
	for (; i < s.length; i++) {
		const c = s[i];
		if (c < '0' || c > '9') {
			break;
		}
		// if (x > 1<<63/10) {
		// 	// overflow
		// 	console.error("x > 1<<63/10")
		// 	return
		// }
		x = x * 10 + (c.charCodeAt(0) - '0'.charCodeAt(0));
		// if (x > 1<<63) {
		// 	// overflow
		// 	console.error("x > 1<<31")
		// 	return
		// }
	}
	return { x, rem: s.slice(i) };
}

function leadingFraction(s: string): { x: number; scale: number; rem: string } {
	let i = 0;
	let scale = 1;
	let x = 0;
	let overflow = false;

	for (; i < s.length; i++) {
		let c = s[i];
		if (c < '0' || c > '9') {
			break;
		}
		if (overflow) {
			continue;
		}
		if (x > (1 << 30) / 10) {
			// It's possible for overflow to give a positive number, so take care.
			overflow = true;
			continue;
		}
		let y = x * 10 + c.charCodeAt(0) - '0'.charCodeAt(0);
		if (y > 1 << 31) {
			overflow = true;
			continue;
		}
		x = y;
		scale *= 10;
	}
	return { x, scale, rem: s.slice(i) };
}
