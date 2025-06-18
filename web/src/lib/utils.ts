import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import type { Duration } from "./types";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, "child"> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any } ? Omit<T, "children"> : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & { ref?: U | null };

export const Nanosecond: Duration = 1
export const Microsecond: Duration = 1000 * Nanosecond
export const Millisecond: Duration = 1000 * Microsecond
export const Second: Duration = 1000 * Millisecond
export const Minute: Duration = 60 * Second
export const Hour: Duration = 60 * Minute

const unitMap = new Map<string, number>([
	["ns", Nanosecond],
	["us", Microsecond],
	["ms", Millisecond],
	["s",  Second],
	["m",  Minute],
	["h",  Hour],
])

// This is a port of golangs time.ParseDuration. 
// Dont ask me how this works, since i just ported it.
export function parseDuration(s: string): Duration | undefined {
	const orig = s
	let d: Duration = 0
	let neg = false

	if (s != "") {
		const c = s[0]
		if (c == '-' || c == '+') {
			neg = c == '-'
			s = s.slice(1)
		}
	}

	if (s == "0") {
		return 0
	}
	if (s == "") {
		console.error(`parseDuration: invalid duration '${orig}'`)
		return undefined
	}


	while (s != "") {
		let v, f: number = 0
		let scale: number = 1 // value = v + f/scale

		// The next character must be [0-9.]
		if (!(s[0] == '.' || '0' <= s[0] && s[0] <= '9')) {
			console.error(`parseDuration: invalid duration '${orig}'`)
			return undefined
		}
		// Consume [0-9]*
		const pl = s.length
		const leading = leadingInt(s)
		if (leading == undefined) {
			console.error(`parseDuration: invalid duration '${orig}'`)
			return undefined
		}
		v = leading.x
		s = leading.rem
		const pre = pl != s.length // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		let post = false
		if (s != "" && s[0] == '.') {
			s = s.slice(1)
			const pl = s.length
			const fract = leadingFraction(s)
			f = fract.x
			scale = fract.scale
			s = fract.rem
			post = pl != s.length
		}
		if (!pre && !post) {
			console.error(`parseDuration: invalid duration '${orig}'`)
			return undefined
		}

		// Consume unit.
		let i = 0
		for (; i < s.length; i++) {
			const c = s[i]
			if (c == '.' || '0' <= c && c <= '9') {
				break
			}
		}
		if (i == 0) {
			console.error(`parseDuration: missing unit in duration '${orig}'`)
			return undefined
		}
		const u = s.slice(0,i)
		s = s.slice(i)
		const unit = unitMap.get(u)
		if (unit == undefined) {
			console.error(`parseDuration: unknown unit '${u}' in duration '${orig}'`)
			return undefined
		}
		if (v > 1<<63/unit) {
			console.error(`parseDuration: invalid duration '${orig}'`)
			return undefined
		}
		v *= unit
		if (f > 0) {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += f * (unit / scale)
			if (v > 1<<31) {
				console.error(`parseDuration: invalid duration '${orig}'`)
				return undefined
			}
		}
		d += v
		if (d > 1<<31) {
			console.error(`parseDuration: invalid duration '${orig}'`)
			return undefined
		}
	}
	
	return d
}

// leadingInt consumes the leading [0-9]* from s.
function leadingInt(s: string): {x: number, rem: string} | undefined {
	let x: number = 0

	let i = 0
	for (; i < s.length; i++) {
		const c = s[i]
		if (c < '0' || c > '9') {
			break
		}
		if (x > 1<<63/10) {
			// overflow
			return 
		}
		x = x*10 + (c.charCodeAt(0) - '0'.charCodeAt(0))
		if (x > 1<<31) {
			// overflow
			return
		}
	}
	return {x, rem: s.slice(i)}
}


function leadingFraction(s: string): {x: number, scale: number, rem: string} {
	let i = 0
	let scale = 1
	let x = 0
	let	overflow = false

	for (; i < s.length; i++) {
		let c = s[i]
		if (c < '0' || c > '9') {
			break
		}
		if (overflow) {
			continue
		}
		if (x > (1<<30)/10) {
			// It's possible for overflow to give a positive number, so take care.
			overflow = true
			continue
		}
		let y = x*10 + c.charCodeAt(0) - '0'.charCodeAt(0)
		if (y > 1<<31) {
			overflow = true
			continue
		}
		x = y
		scale *= 10
	}
	return {x, scale, rem: s.slice(i)}
}
