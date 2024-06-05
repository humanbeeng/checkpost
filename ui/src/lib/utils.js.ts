import { clsx, type ClassValue } from 'clsx';
import { cubicOut } from 'svelte/easing';
import type { TransitionConfig } from 'svelte/transition';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

type FlyAndScaleParams = {
	y?: number;
	x?: number;
	start?: number;
	duration?: number;
};

export const flyAndScale = (
	node: Element,
	params: FlyAndScaleParams = { y: -8, x: 0, start: 0.95, duration: 150 }
): TransitionConfig => {
	const style = getComputedStyle(node);
	const transform = style.transform === 'none' ? '' : style.transform;

	const scaleConversion = (valueA: number, scaleA: [number, number], scaleB: [number, number]) => {
		const [minA, maxA] = scaleA;
		const [minB, maxB] = scaleB;

		const percentage = (valueA - minA) / (maxA - minA);
		const valueB = percentage * (maxB - minB) + minB;

		return valueB;
	};

	const styleToString = (style: Record<string, number | string | undefined>): string => {
		return Object.keys(style).reduce((str, key) => {
			if (style[key] === undefined) return str;
			return str + `${key}:${style[key]};`;
		}, '');
	};

	return {
		duration: params.duration ?? 200,
		delay: 0,
		css: (t) => {
			const y = scaleConversion(t, [0, 1], [params.y ?? 5, 0]);
			const x = scaleConversion(t, [0, 1], [params.x ?? 0, 0]);
			const scale = scaleConversion(t, [0, 1], [params.start ?? 0.95, 1]);

			return styleToString({
				transform: `${transform} translate3d(${x}px, ${y}px, 0) scale(${scale})`,
				opacity: t
			});
		},
		easing: cubicOut
	};
};

export const timeAgo = (timestamp: string) => {
	const now = new Date();
	const then = new Date(timestamp);

	const diffInSeconds = Math.abs(now.getTime() - then.getTime()) / 1000;
	const diffInMinutes = Math.round(diffInSeconds / 60);

	if (diffInMinutes === 0) {
		return 'just now';
	} else if (diffInMinutes < 60) {
		if (diffInMinutes === 1) {
			return `${diffInMinutes} minute ago`;
		}
		return `${diffInMinutes} minutes ago`;
	} else {
		const diffInHours = Math.round(diffInMinutes / 60);
		if (diffInHours === 1) {
			return `${diffInHours} hour ago`;
		}
		return `${diffInHours} hours ago`;
	}
};

export const formatJson = (content: string) => {
	try {
		let obj = JSON.parse(content);
		return JSON.stringify(obj, null, 4);
	} catch (_) {
		return content;
	}
};

export const copy = (content: string) => {
	navigator.clipboard
		.writeText(content)
		.then(() => {
			return true;
		})
		.catch((err: any) => {
			console.error('Failed to copy text: ', err);
		});
};
