import { writable } from 'svelte/store';

type User = {
	name: string;
	email: string;
	plan: Plan;
};

type Plan = 'guest' | 'free' | 'hobby' | 'pro';

export const user = writable<User | null>(null);

export const counter = writable(10);
