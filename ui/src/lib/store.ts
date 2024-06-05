import { writable } from 'svelte/store';
import type { EndpointHistory } from './types';

export const endpointHistory = writable<EndpointHistory>();

export const isFormatEnabled = writable<boolean>(false);

export const selectedRequest = writable<string>('');
