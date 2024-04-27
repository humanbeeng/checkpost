import type { Endpoint } from '@/types';

export type GenerateUrlRequest = {
	endpoint: string;
};

export type GenerateUrlResponse = {
	url: string;
	expires_at: string;
	plan: Plan;
};

export type UserEndpointsResponse = {
	endpoints: Endpoint[];
};
