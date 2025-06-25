import { get, writable, type Writable } from 'svelte/store';

export function localStore<T = unknown>(key: string, data: T): Writable<T> {
	const store = writable(data);
	const { subscribe, set } = store;
	const isBrowser = typeof window !== 'undefined';

	isBrowser && localStorage[key] && set(JSON.parse(localStorage[key]));

	isBrowser && !localStorage[key] && (localStorage[key] = JSON.stringify(data));

	return {
		subscribe,
		set: (n) => {
			isBrowser && (localStorage[key] = JSON.stringify(n));
			set(n);
		},
		update: (cb) => {
			const updatedStore = cb(get(store));

			isBrowser && (localStorage[key] = JSON.stringify(updatedStore));
			set(updatedStore);
		}
	};
}
