import { writable } from 'svelte/store';
import type { UrlHistory } from './types';

export const urlHistory = writable<UrlHistory>();

export const isFormatEnabled = writable<boolean>(false);

export const selectedRequest = writable<string>('');
