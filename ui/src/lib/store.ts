import { writable } from 'svelte/store';
import type { User } from './types';

type Plan = 'guest' | 'free' | 'hobby' | 'pro';

export const user = writable<User | null>(null);
