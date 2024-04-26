export function debounce(fn: any, delay: number) {
	let timeoutId: number;
	return (...args: any[]) => {
		clearTimeout(timeoutId);
		timeoutId = setTimeout(() => {
			fn.apply(null, args);
		}, delay);
	};
}
